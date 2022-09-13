package util

const (
	AP_NORTHEAST_2 string = "ap-northeast-2" // 서울
	EU_WEST_1      string = "eu-west-1"      // 아일랜드
)

var AwsRegionSet []string = []string{
	AP_NORTHEAST_2,
	EU_WEST_1,
}

const (
	UBUNTU string = "ubuntu"
	RHEL   string = "rhel"
)

var AwsSupportOs []string = []string{
	UBUNTU,
	RHEL,
}

type AwsParams struct {
	AmiName    string
	AmiOwner   string
	User       string
	SshKeyName string
}
