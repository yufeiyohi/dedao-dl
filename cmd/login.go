package cmd

import (
	"os"
	"strconv"
	"strings"

	"errors"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
)

// Cookie cookie from https://www.dedao.cn
var Cookie string
var qr bool

// Login login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录得到 pc 端 https://www.dedao.cn",
	Long:  `使用 dedao-dl login to login https://www.dedao.cn`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 确保配置初始化
		err := config.Instance.Init()
		if err != nil && err.Error() != "未登陆" && !strings.Contains(err.Error(), "存在登录的用户") {
			return err
		}

		if qr {
			err := app.LoginByQr()
			return err
		}
		if Cookie == "" {
			// defaultCookie := app.GetCookie()
			// if defaultCookie == "" {
			// 	return errors.New("自动获取 cookie 失败，请使用参数设置 cookie")
			// }
			// Cookie = defaultCookie
			return errors.New("请使用参数设置 cookie")
		} else {
			err := app.LoginByCookie(Cookie)
			return err
		}
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		who()
	},
}

var whoCmd = &cobra.Command{
	Use:     "who",
	Short:   "查看当前登录的用户",
	Long:    `使用 dedao-dl who 当前登录的用户信息`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		who()
	},
}

var usersCmd = &cobra.Command{
	Use:     "users",
	Short:   "查看登录过的用户列表",
	Long:    `使用 dedao-dl users 查看登录过的用户列表`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		users()
	},
}

var suCmd = &cobra.Command{
	Use:     "su",
	Short:   "切换登录过的账号",
	Long:    `使用 dedao-dl su 切换当前登录的账号`,
	Args:    cobra.ExactArgs(1),
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("please input UID")
		}
		uid := args[0]
		err := app.SwitchAccount(uid)
		return err
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		who()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(suCmd)
	loginCmd.Flags().StringVarP(&Cookie, "cookie", "c", "", "cookie from https://www.dedao.cn")
	loginCmd.Flags().BoolVarP(&qr, "qrcode", "q", false, "scan qrcode login from https://www.dedao.cn")
}

// current login user
func who() {
	activeUser := config.Instance.ActiveUser()
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"UID", "姓名", "头像"})
	table.Append([]string{activeUser.UIDHazy, activeUser.Name, activeUser.Avatar})
	table.Render()
}

// users get login user list
func users() {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "UID", "姓名", "头像"})
	for i, user := range config.Instance.Users {
		table.Append([]string{strconv.Itoa(i), user.UIDHazy, user.Name, user.Avatar})
	}
	table.Render()

}
