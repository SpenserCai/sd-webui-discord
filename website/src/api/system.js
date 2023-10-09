/*
 * @Author: SpenserCai
 * @Date: 2023-10-09 11:05:23
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 17:52:42
 * @Description: file content
 */
import request from "@/utils/request";

export function discordserver() {
  return request({
    url: "/discord_server",
    method: "get"
  });
}

export function cluster() {
  return request({
    url: "/cluster",
    method: "get"
  });
}