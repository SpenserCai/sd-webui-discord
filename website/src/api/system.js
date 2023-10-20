/*
 * @Author: SpenserCai
 * @Date: 2023-10-09 11:05:23
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-20 13:51:14
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

export function userlist(body) {
  return request({
    url: "/user_list",
    method: "post",
    data: body
  });
}

export function setuserprivate(body) {
  return request({
    url: "/set_user_private",
    method: "post",
    data: body
  });
}