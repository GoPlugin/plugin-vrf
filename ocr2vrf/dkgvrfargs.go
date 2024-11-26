package ocr2vrf

import (
	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	dkg_contract "github.com/goplugin/plugin-vrf/internal/dkg/contract"
	vrf_types "github.com/goplugin/plugin-vrf/types"
)

type DKGVRFArgs struct {
	DKGLogger commontypes.Logger
	VRFLogger commontypes.Logger

	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	V2Bootstrappers []commontypes.BootstrapperLocator

	OffchainKeyring types.OffchainKeyring

	OnchainKeyring types.OnchainKeyring

	DKGOffchainConfigDigester types.OffchainConfigDigester

	VRFOffchainConfigDigester types.OffchainConfigDigester

	VRFContractConfigTracker types.ContractConfigTracker

	VRFContractTransmitter types.ContractTransmitter

	VRFDatabase types.Database

	VRFLocalConfig types.LocalConfig

	VRFMonitoringEndpoint commontypes.MonitoringEndpoint

	DKGContractConfigTracker types.ContractConfigTracker

	DKGContract dkg_contract.OnchainContract

	DKGContractTransmitter types.ContractTransmitter

	DKGDatabase types.Database

	DKGLocalConfig types.LocalConfig

	DKGMonitoringEndpoint commontypes.MonitoringEndpoint

	DKGSharePersistence vrf_types.DKGSharePersistence

	Serializer         vrf_types.ReportSerializer
	JuelsPerFeeCoin    vrf_types.JuelsPerFeeCoin
	ReasonableGasPrice vrf_types.ReasonableGasPrice
	Coordinator        vrf_types.CoordinatorInterface

	ConfirmationDelays []uint32

	Esk   dkg_contract.EncryptionSecretKey
	Ssk   dkg_contract.SigningSecretKey
	KeyID dkg_contract.KeyID

	DKGReportingPluginFactoryDecorator func(factory types.ReportingPluginFactory) types.ReportingPluginFactory
	VRFReportingPluginFactoryDecorator func(factory types.ReportingPluginFactory) types.ReportingPluginFactory
}
