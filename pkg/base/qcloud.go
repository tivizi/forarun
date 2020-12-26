package base

import (
	"fmt"
	"log"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"github.com/tivizi/forarun/approot/config"
)

var cdnCli *cdn.Client

func init() {
	qcloud := config.GetContext().QcloudConfig
	if qcloud.EnableCDN {
		credential := common.NewCredential(qcloud.SecretID, qcloud.SecretKey)
		cpf := profile.NewClientProfile()
		cli, err := cdn.NewClient(credential, regions.Guangzhou, cpf)
		if err != nil {
			panic(err)
		}
		cdnCli = cli
		log.Println("Qcloud CDN: Enabled")
	}

}

func addDomain(cli *cdn.Client, domain, origin string) {
	serviceType := "web"
	originType := "domain"
	area := "overseas"
	req := cdn.NewAddCdnDomainRequest()
	req.Domain = &domain
	req.ServiceType = &serviceType
	req.Origin = &cdn.Origin{
		Origins:    []*string{&origin},
		OriginType: &originType,
	}
	req.Area = &area
	res, err := cli.AddCdnDomain(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToJsonString())
}

func stopDomain(cli *cdn.Client, domain string) {
	req := cdn.NewStopCdnDomainRequest()
	req.Domain = &domain
	res, err := cli.StopCdnDomain(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToJsonString())
}

func delDomain(cli *cdn.Client, domain string) {
	req := cdn.NewDeleteCdnDomainRequest()
	req.Domain = &domain
	res, err := cli.DeleteCdnDomain(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToJsonString())
}

func getDomain(cli *cdn.Client, domain string) string {
	key := "domain"
	req := cdn.NewDescribeDomainsConfigRequest()
	req.Filters = append(req.Filters, &cdn.DomainFilter{
		Name:  &key,
		Value: []*string{&domain},
	})
	res, err := cli.DescribeDomainsConfig(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToJsonString())
	return *(res.Response.Domains[0].Cname)
}

func updateDomain(cli *cdn.Client, domain, certID string) {
	onOff := "on"
	redirectType := "https"
	redirectStatusCode := int64(302)
	req := cdn.NewUpdateDomainConfigRequest()
	req.Domain = &domain
	req.Https = &cdn.Https{
		Switch: &onOff,
		CertInfo: &cdn.ServerCert{
			CertId: &certID,
		},
	}
	req.ForceRedirect = &cdn.ForceRedirect{
		Switch:             &onOff,
		RedirectType:       &redirectType,
		RedirectStatusCode: &redirectStatusCode,
	}
	res, err := cli.UpdateDomainConfig(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToJsonString())
}
