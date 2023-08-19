'''
Author: SpenserCai
Date: 2023-08-19 15:41:07
version: 
LastEditors: SpenserCai
LastEditTime: 2023-08-19 20:55:25
Description: file content
'''

import argparse
import os


temp='''
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) {Cmd}Options() *discordgo.ApplicationCommand {{
	return &discordgo.ApplicationCommand{{
		Name:        "{cmd}",
		Description: "Remove background from image",
		Options: []*discordgo.ApplicationCommandOption{{
        
		}},
	}}
}}

func (shdl SlashHandler) {Cmd}SetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.{RequestName}Request) {{

	for _, v := range dsOpt {{
		switch v.Name {{
        
		}}
	}}
}}

func (shdl SlashHandler) {Cmd}Action(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.{RequestName}Request, node *cluster.ClusterNode) {{
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {{
		log.Println(err)
		return
	}}
	{cmd} := &intersvc.{RequestName}{{RequestItem: opt}}
	{cmd}.Action(node.StableClient)
	if {cmd}.Error != nil {{
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{{
			Content: func() *string {{ v := {cmd}.Error.Error(); return &v }}(),
		}})
	}}
}}

func (shdl SlashHandler) {Cmd}CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {{
	option := &intersvc.{RequestName}Request{{}}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{{}}, error) {{
		shdl.{Cmd}SetOptions(i.ApplicationCommandData().Options, option)
		shdl.{Cmd}Action(s, i, option, node)
		return nil, nil
	}}
	callback := func() {{}}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}}
'''

parser = argparse.ArgumentParser(description='Generate command template')
parser.add_argument('--cmd', type=str, help='command name',required=True)
parser.add_argument('--rn', type=str, help='request name',required=True)
args = parser.parse_args()



# cmd是cmd的参数，Cmd是cmd用_分割后的首字母大写
cmd = args.cmd
Cmd = ''.join([i.capitalize() for i in cmd.split('_')])
rn = args.rn

outCode = temp.format(cmd=cmd, Cmd=Cmd, RequestName=rn)

# 写入次文件所在目录的../../bot/slash_handler/{cmd}.go中如果文件已经存在则不写入

# 获取当前文件的所在路径
curPath = os.path.abspath(__file__)
outPath = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(curPath))), 'dbot', 'slash_handler', f'{cmd}.go')

if os.path.exists(outPath):
    print(f'{outPath} already exists')
else:
    with open(outPath, 'w') as f:
        f.write(outCode)
    print(f'{outPath} generated')
