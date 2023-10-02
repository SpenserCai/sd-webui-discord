/*
 * @Author: SpenserCai
 * @Date: 2023-10-02 21:33:34
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-02 21:52:46
 * @Description: file content
 */
import request from "@/utils/request";

export function userinfo() {
  return request({
    url: "/user_info",
    method: "get"
  });
}