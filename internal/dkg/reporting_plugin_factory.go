package dkg

import (
	"context"
	"fmt"
	"sync"

	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-vrf/internal/crypto/player_idx"
	"github.com/goplugin/plugin-vrf/internal/util"
)

type dkgReportingPluginFactory struct {
	l    *localArgs
	lock sync.RWMutex

	dkgInProgress bool

	testmode          bool
	xxxDKGTestingOnly *dkg
}

var _ types.ReportingPluginFactory = (*dkgReportingPluginFactory)(nil)

func (d *dkgReportingPluginFactory) NewReportingPlugin(
	c types.ReportingPluginConfig,
) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	emptyInfo := types.ReportingPluginInfo{}
	if d.dkgInProgress {
		return nil, emptyInfo, fmt.Errorf(
			"attempt to initiate DKG round while an earlier DKG round is in progress",
		)
	}
	d.dkgInProgress = true
	a, err := unmarshalPluginConfig(c.OffchainConfig, c.OnchainConfig)
	if err != nil {
		return nil, emptyInfo,
			util.WrapError(err, "could not read offchain plugin config")
	}
	if c.N > int(player_idx.MaxPlayer) {
		return nil, emptyInfo,
			fmt.Errorf("too many players: %d > %d", c.N, player_idx.MaxPlayer)
	}
	args, err := a.NewDKGArgs(
		c.ConfigDigest, d.l, c.OracleID, player_idx.Int(c.N), player_idx.Int(c.F),
	)
	if err != nil {
		return nil, emptyInfo, util.WrapError(err, "could not construct DKG args")
	}
	d.l.logger.Debug("constructing share set", commontypes.LogFields{})
	dkg, err := d.NewDKG(args)
	if err != nil {
		return nil, emptyInfo, util.WrapError(err, "while creating reporting plugin")
	}
	d.l.logger.Debug("finished constructing share set", commontypes.LogFields{})
	if d.testmode {
		d.xxxDKGTestingOnly = dkg
	}
	dkg.keyConsumer.KeyInvalidated(dkg.keyID)
	return dkg, types.ReportingPluginInfo{
		Name: fmt.Sprintf("dkg instance %v", dkg.selfIdx),
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       1000,
			MaxObservationLength: 1_000_000,
			MaxReportLength:      10_000,
		},
		UniqueReports: true,
	}, nil
}

func (d *dkgReportingPluginFactory) NewDKG(a *NewDKGArgs) (*dkg, error) {
	if err := a.SanityCheckArgs(); err != nil {
		return nil, util.WrapError(err, "could not construct new DKG")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	factory := &dkg{
		a.t,
		sync.RWMutex{},
		a.selfIdx,
		a.cfgDgst,
		a.keyID,
		a.keyConsumer,
		newShareRecords(),
		nil,
		a.esk,
		a.epks,
		a.ssk,
		a.spks,
		a.encryptionGroup,
		a.translationGroup,
		a.translator,
		nil,
		a.contract,
		false,
		d.markCompleted,
		a.db,
		a.logger,
		a.randomness,
		ctx,
		cancelFunc,
	}
	if err := factory.initializeShareSets(a.signingGroup()); err != nil {
		return nil, util.WrapError(err, "could not initialize share sets")
	}

	res := make(chan error, 1)
	go func(ctx context.Context) {
		if factory.keyReportedOnchain(ctx) {
			err := factory.recoverDistributedKeyShare(ctx)
			if err != nil {
				errMsg := "could not reconstruct shares for an available distributed key"
				err2 := util.WrapError(err, errMsg)
				res <- err2
			}
		}
		res <- nil
	}(ctx)
	err := <-res
	if err != nil {
		return nil, err
	}
	return factory, nil
}

func (d *dkgReportingPluginFactory) markCompleted() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.dkgInProgress = false
}

func (d *dkgReportingPluginFactory) SetKeyConsumer(k KeyConsumer) {
	d.l.keyConsumer = k
}
