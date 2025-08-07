package machinery

// RegisterTask 注册任务处理函数到Machinery服务器
// 将任务映射表中的所有任务处理函数注册到服务器，使其可以被远程调用
// 参数:
//   - server: Machinery服务器实例
//   - taskMap: 任务名称到处理函数的映射
func (m *Machinery) RegisterTask() {
	m.server.RegisterTasks(map[string]any{
		"get_access_token":                   m.oceanengine.GetAccessToken,
		"refresh_access_token":               m.oceanengine.RefreshAccessToken,
		"oceanengine_get_advertiser_id_list": m.oceanengine.GetAdvertiserIdList,
		"oceanengine_get_report_advertiser":  m.oceanengine.ReportAdvertiserGet,
	})
}
