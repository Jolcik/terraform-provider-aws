// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package names provides constants for AWS service names that are used as keys
// for the endpoints slice in internal/conns/conns.go. The package also exposes
// access to data found in the data/names_data.csv file, which provides additional
// service-related name information.
//
// Consumers of the names package include the conns package
// (internal/conn/conns.go), the provider package
// (internal/provider/provider.go), generators, and the skaff tool.
//
// It is very important that information in the data/names_data.csv be exactly
// correct because the Terrform AWS Provider relies on the information to
// function correctly.
package names

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-aws/names/data"
	"golang.org/x/exp/slices"
)

// These "should" be defined by the AWS Go SDK v2, but currently aren't.
const (
	AccessAnalyzerEndpointID             = "access-analyzer"
	AccountEndpointID                    = "account"
	ACMEndpointID                        = "acm"
	AMPEndpointID                        = "aps"
	AppFlowEndpointID                    = "appflow"
	AppRunnerEndpointID                  = "apprunner"
	AthenaEndpointID                     = "athena"
	AuditManagerEndpointID               = "auditmanager"
	BedrockEndpointID                    = "bedrock"
	ChimeSDKVoiceEndpointID              = "voice-chime"
	ChimeSDKMediaPipelinesEndpointID     = "media-pipelines-chime"
	CleanRoomsEndpointID                 = "cleanrooms"
	CloudWatchLogsEndpointID             = "logs"
	CodeDeployEndpointID                 = "codedeploy"
	CodeGuruProfilerEndpointID           = "codeguru-profiler"
	CodeStarConnectionsEndpointID        = "codestar-connections"
	CodeStarNotificationsEndpointID      = "codestar-notifications"
	ComprehendEndpointID                 = "comprehend"
	ComputeOptimizerEndpointID           = "computeoptimizer"
	DocDBElasticEndpointID               = "docdb-elastic"
	ControlTowerEndpointID               = "controltower"
	DSEndpointID                         = "ds"
	ECREndpointID                        = "api.ecr"
	EKSEndpointID                        = "eks"
	EMREndpointID                        = "elasticmapreduce"
	EMRServerlessEndpointID              = "emrserverless"
	EvidentlyEndpointID                  = "evidently"
	GlacierEndpointID                    = "glacier"
	GroundStationEndpointID              = "groundstation"
	IdentityStoreEndpointID              = "identitystore"
	Inspector2EndpointID                 = "inspector2"
	InternetMonitorEndpointID            = "internetmonitor"
	IVSChatEndpointID                    = "ivschat"
	KendraEndpointID                     = "kendra"
	KeyspacesEndpointID                  = "keyspaces"
	LambdaEndpointID                     = "lambda"
	LexV2ModelsEndpointID                = "models-v2-lex"
	MediaLiveEndpointID                  = "medialive"
	ObservabilityAccessManagerEndpointID = "oam"
	OpenSearchServerlessEndpointID       = "aoss"
	PipesEndpointID                      = "pipes"
	PollyEndpointID                      = "polly"
	PricingEndpointID                    = "pricing"
	QLDBEndpointID                       = "qldb"
	RedshiftDataEndpointID               = "redshift-data"
	ResourceExplorer2EndpointID          = "resource-explorer-2"
	ResourceGroupsEndpointID             = "resource-groups"
	ResourceGroupsTaggingAPIEndpointID   = "tagging"
	RolesAnywhereEndpointID              = "rolesanywhere"
	Route53DomainsEndpointID             = "route53domains"
	SchedulerEndpointID                  = "scheduler"
	SecurityLakeEndpointID               = "securitylake"
	ServiceQuotasEndpointID              = "servicequotas"
	S3EndpointID                         = "s3"
	S3ControlEndpointID                  = "s3-control"
	SecurityHubEndpointID                = "securityhub"
	SESV2EndpointID                      = "sesv2"
	SNSEndpointID                        = "sns"
	SQSEndpointID                        = "sqs"
	SSMEndpointID                        = "ssm"
	SSMContactsEndpointID                = "ssm-contacts"
	SSMIncidentsEndpointID               = "ssm-incidents"
	SSOAdminEndpointID                   = "sso"
	STSEndpointID                        = "sts"
	SWFEndpointID                        = "swf"
	TimestreamWriteEndpointID            = "ingest.timestream"
	TranscribeEndpointID                 = "transcribe"
	VerifiedPermissionsEndpointID        = "verifiedpermissions"
	VPCLatticeEndpointID                 = "vpc-lattice"
	XRayEndpointID                       = "xray"
)

// These should move to aws-sdk-go-base.
// See https://github.com/hashicorp/aws-sdk-go-base/issues/649.
const (
	ChinaPartitionID      = "aws-cn"     // AWS China partition.
	ISOPartitionID        = "aws-iso"    // AWS ISO (US) partition.
	ISOBPartitionID       = "aws-iso-b"  // AWS ISOB (US) partition.
	ISOEPartitionID       = "aws-iso-e"  // AWS ISOE (Europe) partition.
	ISOFPartitionID       = "aws-iso-f"  // AWS ISOF partition.
	StandardPartitionID   = "aws"        // AWS Standard partition.
	USGovCloudPartitionID = "aws-us-gov" // AWS GovCloud (US) partition.
)

const (
	// AWS Standard partition's regions.
	GlobalRegionID = "aws-global" // AWS Standard global region.

	AFSouth1RegionID     = "af-south-1"     // Africa (Cape Town).
	APEast1RegionID      = "ap-east-1"      // Asia Pacific (Hong Kong).
	APNortheast1RegionID = "ap-northeast-1" // Asia Pacific (Tokyo).
	APNortheast2RegionID = "ap-northeast-2" // Asia Pacific (Seoul).
	APNortheast3RegionID = "ap-northeast-3" // Asia Pacific (Osaka).
	APSouth1RegionID     = "ap-south-1"     // Asia Pacific (Mumbai).
	APSouth2RegionID     = "ap-south-2"     // Asia Pacific (Hyderabad).
	APSoutheast1RegionID = "ap-southeast-1" // Asia Pacific (Singapore).
	APSoutheast2RegionID = "ap-southeast-2" // Asia Pacific (Sydney).
	APSoutheast3RegionID = "ap-southeast-3" // Asia Pacific (Jakarta).
	APSoutheast4RegionID = "ap-southeast-4" // Asia Pacific (Melbourne).
	CACentral1RegionID   = "ca-central-1"   // Canada (Central).
	CAWest1RegionID      = "ca-west-1"      // Canada West (Calgary).
	EUCentral1RegionID   = "eu-central-1"   // Europe (Frankfurt).
	EUCentral2RegionID   = "eu-central-2"   // Europe (Zurich).
	EUNorth1RegionID     = "eu-north-1"     // Europe (Stockholm).
	EUSouth1RegionID     = "eu-south-1"     // Europe (Milan).
	EUSouth2RegionID     = "eu-south-2"     // Europe (Spain).
	EUWest1RegionID      = "eu-west-1"      // Europe (Ireland).
	EUWest2RegionID      = "eu-west-2"      // Europe (London).
	EUWest3RegionID      = "eu-west-3"      // Europe (Paris).
	ILCentral1RegionID   = "il-central-1"   // Israel (Tel Aviv).
	MECentral1RegionID   = "me-central-1"   // Middle East (UAE).
	MESouth1RegionID     = "me-south-1"     // Middle East (Bahrain).
	SAEast1RegionID      = "sa-east-1"      // South America (Sao Paulo).
	USEast1RegionID      = "us-east-1"      // US East (N. Virginia).
	USEast2RegionID      = "us-east-2"      // US East (Ohio).
	USWest1RegionID      = "us-west-1"      // US West (N. California).
	USWest2RegionID      = "us-west-2"      // US West (Oregon).

	// AWS China partition's regions.
	CNNorth1RegionID     = "cn-north-1"     // China (Beijing).
	CNNorthwest1RegionID = "cn-northwest-1" // China (Ningxia).

	// AWS GovCloud (US) partition's regions.
	USGovEast1RegionID = "us-gov-east-1" // AWS GovCloud (US-East).
	USGovWest1RegionID = "us-gov-west-1" // AWS GovCloud (US-West).

	// AWS ISO (US) partition's regions.
	USISOEast1RegionID = "us-iso-east-1" // US ISO East.
	USISOWest1RegionID = "us-iso-west-1" // US ISO WEST.

	// AWS ISOB (US) partition's regions.
	USISOBEast1RegionID = "us-isob-east-1" // US ISOB East (Ohio).
)

func DNSSuffixForPartition(partition string) string {
	switch partition {
	case "":
		return ""
	case ChinaPartitionID:
		return "amazonaws.com.cn"
	case ISOPartitionID:
		return "c2s.ic.gov"
	case ISOBPartitionID:
		return "sc2s.sgov.gov"
	case ISOEPartitionID:
		return "cloud.adc-e.uk"
	case ISOFPartitionID:
		return "csp.hci.ic.gov"
	default:
		return "amazonaws.com"
	}
}

func PartitionForRegion(region string) string {
	switch region {
	case "":
		return ""
	case CNNorth1RegionID, CNNorthwest1RegionID:
		return ChinaPartitionID
	case USISOEast1RegionID, USISOWest1RegionID:
		return ISOPartitionID
	case USISOBEast1RegionID:
		return ISOBPartitionID
	case USGovEast1RegionID, USGovWest1RegionID:
		return USGovCloudPartitionID
	default:
		return StandardPartitionID
	}
}

// ReverseDNS switches a DNS hostname to reverse DNS and vice-versa.
func ReverseDNS(hostname string) string {
	parts := strings.Split(hostname, ".")

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}

// Type ServiceDatum corresponds closely to columns in `data/names_data.csv` and are
// described in detail in README.md.
type ServiceDatum struct {
	Aliases            []string
	Brand              string
	DeprecatedEnvVar   string
	EndpointOnly       bool
	GoV1ClientTypeName string
	GoV1Package        string
	GoV2Package        string
	HumanFriendly      string
	ProviderNameUpper  string
	TfAwsEnvVar        string
}

// serviceData key is the AWS provider service package
var serviceData map[string]*ServiceDatum

func init() {
	serviceData = make(map[string]*ServiceDatum)

	// Data from names_data.csv
	if err := readCSVIntoServiceData(); err != nil {
		log.Fatalf("reading CSV into service data: %s", err)
	}
}

func readCSVIntoServiceData() error {
	// names_data.csv is dynamically embedded so changes, additions should be made
	// there also

	d, err := data.ReadAllServiceData()
	if err != nil {
		return fmt.Errorf("reading CSV into service data: %w", err)
	}

	for _, l := range d {
		if l.Exclude() {
			continue
		}

		if l.NotImplemented() && !l.EndpointOnly() {
			continue
		}

		p := l.ProviderPackage()

		serviceData[p] = &ServiceDatum{
			Brand:              l.Brand(),
			DeprecatedEnvVar:   l.DeprecatedEnvVar(),
			EndpointOnly:       l.EndpointOnly(),
			GoV1ClientTypeName: l.GoV1ClientTypeName(),
			GoV1Package:        l.GoV1Package(),
			GoV2Package:        l.GoV2Package(),
			HumanFriendly:      l.HumanFriendly(),
			ProviderNameUpper:  l.ProviderNameUpper(),
			TfAwsEnvVar:        l.TfAwsEnvVar(),
		}

		a := []string{p}

		if len(l.Aliases()) > 0 {
			a = append(a, l.Aliases()...)
		}

		serviceData[p].Aliases = a
	}

	return nil
}

func ProviderPackageForAlias(serviceAlias string) (string, error) {
	for k, v := range serviceData {
		for _, hclKey := range v.Aliases {
			if serviceAlias == hclKey {
				return k, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find service for service alias %s", serviceAlias)
}

func ProviderPackages() []string {
	keys := make([]string, len(serviceData))

	i := 0
	for k := range serviceData {
		keys[i] = k
		i++
	}

	return keys
}

func Aliases() []string {
	keys := make([]string, 0)

	for _, v := range serviceData {
		keys = append(keys, v.Aliases...)
	}

	return keys
}

type Endpoint struct {
	ProviderPackage string
	Aliases         []string
}

func Endpoints() []Endpoint {
	endpoints := make([]Endpoint, 0, len(serviceData))

	for k, v := range serviceData {
		ep := Endpoint{
			ProviderPackage: k,
		}
		if len(v.Aliases) > 1 {
			idx := slices.Index(v.Aliases, k)
			if idx != -1 {
				aliases := slices.Delete(v.Aliases, idx, idx+1)
				ep.Aliases = aliases
			}
		}
		endpoints = append(endpoints, ep)
	}

	return endpoints
}

type ServiceNameUpper struct {
	ProviderPackage   string
	ProviderNameUpper string
}

func ServiceNamesUpper() []ServiceNameUpper {
	serviceNames := make([]ServiceNameUpper, 0, len(serviceData))

	for k, v := range serviceData {
		sn := ServiceNameUpper{
			ProviderPackage:   k,
			ProviderNameUpper: v.ProviderNameUpper,
		}
		serviceNames = append(serviceNames, sn)
	}

	return serviceNames
}

func ProviderNameUpper(service string) (string, error) {
	if v, ok := serviceData[service]; ok {
		return v.ProviderNameUpper, nil
	}

	return "", fmt.Errorf("no service data found for %s", service)
}

func DeprecatedEnvVar(service string) string {
	if v, ok := serviceData[service]; ok {
		return v.DeprecatedEnvVar
	}

	return ""
}

func TfAwsEnvVar(service string) string {
	if v, ok := serviceData[service]; ok {
		return v.TfAwsEnvVar
	}

	return ""
}

func FullHumanFriendly(service string) (string, error) {
	if v, ok := serviceData[service]; ok {
		if v.Brand == "" {
			return v.HumanFriendly, nil
		}

		return fmt.Sprintf("%s %s", v.Brand, v.HumanFriendly), nil
	}

	if s, err := ProviderPackageForAlias(service); err == nil {
		return FullHumanFriendly(s)
	}

	return "", fmt.Errorf("no service data found for %s", service)
}

func HumanFriendly(service string) (string, error) {
	if v, ok := serviceData[service]; ok {
		return v.HumanFriendly, nil
	}

	if s, err := ProviderPackageForAlias(service); err == nil {
		return HumanFriendly(s)
	}

	return "", fmt.Errorf("no service data found for %s", service)
}

func AWSGoPackage(providerPackage string, version int) (string, error) {
	switch version {
	case 1:
		return AWSGoV1Package(providerPackage)
	case 2:
		return AWSGoV2Package(providerPackage)
	default:
		return "", fmt.Errorf("unsupported AWS SDK Go version: %d", version)
	}
}

func AWSGoV1Package(providerPackage string) (string, error) {
	if v, ok := serviceData[providerPackage]; ok {
		return v.GoV1Package, nil
	}

	return "", fmt.Errorf("getting AWS SDK Go v1 package, %s not found", providerPackage)
}

func AWSGoV2Package(providerPackage string) (string, error) {
	if v, ok := serviceData[providerPackage]; ok {
		return v.GoV2Package, nil
	}

	return "", fmt.Errorf("getting AWS SDK Go v2 package, %s not found", providerPackage)
}

func AWSGoClientTypeName(providerPackage string, version int) (string, error) {
	switch version {
	case 1:
		return AWSGoV1ClientTypeName(providerPackage)
	case 2:
		return "Client", nil
	default:
		return "", fmt.Errorf("unsupported AWS SDK Go version: %d", version)
	}
}

func AWSGoV1ClientTypeName(providerPackage string) (string, error) {
	if v, ok := serviceData[providerPackage]; ok {
		return v.GoV1ClientTypeName, nil
	}

	return "", fmt.Errorf("getting AWS SDK Go v1 client type name, %s not found", providerPackage)
}
