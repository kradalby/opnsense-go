package opnsense

func (c *Client) PowerOff() (*StatusMessage, error) {
	api := "core/firmware/poweroff"

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) Reboot() (*StatusMessage, error) {
	api := "core/firmware/reboot"

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) Upgrade() (*StatusMessage, error) {
	api := "core/firmware/upgrade"

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

type UpgradeStatusMessage struct {
	Status string `json:"status"`
	Log    string `json:"log"`
}

func (c *Client) UpgradeStatus() (*UpgradeStatusMessage, error) {
	api := "core/firmware/upgradestatus"

	var status UpgradeStatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) Audit() (*StatusMessage, error) {
	api := "core/firmware/audit"

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

type FirmwareConfig struct {
	Flavour string `json:"flavour"`
	Mirror  string `json:"mirror"`
	Type    string `json:"type"`
}

func (c *Client) GetFirmwareConfig() (*FirmwareConfig, error) {
	api := "core/firmware/getfirmwareconfig"

	var firmwareConfig FirmwareConfig
	err := c.GetAndUnmarshal(api, &firmwareConfig)

	return &firmwareConfig, err
}

func (c *Client) SetFirmwareConfig(config FirmwareConfig) (*StatusMessage, error) {
	api := "core/firmware/setfirmwareconfig"

	var status StatusMessage

	err := c.PostAndMarshal(api, config, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

type FirmwareOptions struct {
	HasSubscription []string            `json:"has_subscription"`
	Flavours        []map[string]string `json:"flavours"`
	Families        []map[string]string `json:"families"`
	Mirrors         []map[string]string `json:"mirrors"`
	AllowCustom     bool                `json:"allow_custom"`
}

func (c *Client) GetFirmwareOptions() (*FirmwareOptions, error) {
	api := "core/firmware/getfirmwareoptions"

	var firmwareOptions FirmwareOptions
	err := c.GetAndUnmarshal(api, &firmwareOptions)

	return &firmwareOptions, err
}

type Status struct {
	Connection          string   `json:"connection"`
	DowngradePackages   []string `json:"downgrade_packages"`
	DownloadSize        string   `json:"download_size"`
	LastCheck           string   `json:"last_check"`
	NewPackages         []string `json:"new_packages"`
	OsVersion           string   `json:"os_version"`
	ProductName         string   `json:"product_name"`
	ProductVersion      string   `json:"product_version"`
	ReinstallPackages   []string `json:"reinstall_packages"`
	RemovePackages      []string `json:"remove_packages"`
	Repository          string   `json:"repository"`
	Updates             string   `json:"updates"`
	UpgradeMajorMessage string   `json:"upgrade_major_message"`
	UpgradeMajorVersion string   `json:"upgrade_major_version"`
	UpgradeNeedsReboot  string   `json:"upgrade_needs_reboot"`
	UpgradePackages     []string `json:"upgrade_packages"`
	AllPackages         []string `json:"all_packages"`
	StatusMsg           string   `json:"status_msg"`
	Status              string   `json:"status"`
}

func (c *Client) GetStatus() (*Status, error) {
	api := "core/firmware/status"

	var status Status
	err := c.GetAndUnmarshal(api, &status)

	return &status, err
}

type UpgradeStatus struct {
	Status string `json:"status"`
	Log    string `json:"log"`
}

func (c *Client) GetUpgradeStatus() (*UpgradeStatus, error) {
	api := "core/firmware/upgradestatus"

	var status UpgradeStatus
	err := c.GetAndUnmarshal(api, &status)

	return &status, err
}

type ChangeLog struct {
	Series  string `json:"series"`
	Version string `json:"version"`
	Date    string `json:"date"`
}

type Information struct {
	ProductVersion string      `json:"product_version"`
	ProductName    string      `json:"product_name"`
	Package        []Package   `json:"package"`
	Plugin         []Package   `json:"plugin"`
	Changelog      []ChangeLog `json:"changelog"`
}

func (c *Client) GetInformation() (*Information, error) {
	api := "core/firmware/info"

	var information Information
	err := c.GetAndUnmarshal(api, &information)

	return &information, err
}

type Package struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Comment    string `json:"comment"`
	Flatsize   string `json:"flatsize"`
	Locked     string `json:"locked"`
	License    string `json:"license"`
	Repository string `json:"repository"`
	Origin     string `json:"origin"`
	Provided   Bool   `json:"provided"`
	Installed  Bool   `json:"installed"`
	Path       string `json:"path"`
	Configured Bool   `json:"configured"`
}

func (c *Client) PackageInstall(packageName string) (*StatusMessage, error) {
	api := "core/firmware/install/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageReInstall(packageName string) (*StatusMessage, error) {
	api := "core/firmware/reinstall/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageRemove(packageName string) (*StatusMessage, error) {
	api := "core/firmware/remove/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageLock(packageName string) (*StatusMessage, error) {
	api := "core/firmware/lock/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageUnlock(packageName string) (*StatusMessage, error) {
	api := "core/firmware/unlock/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageDetails(packageName string) (*StatusMessage, error) {
	api := "core/firmware/details/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) PackageLicense(packageName string) (*StatusMessage, error) {
	api := "core/firmware/license/" + packageName

	var status StatusMessage

	err := c.PostAndMarshal(api, nil, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}
