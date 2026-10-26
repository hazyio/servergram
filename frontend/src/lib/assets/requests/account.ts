import { BaseRequest } from "."
export class AccountRequests extends BaseRequest {
  base_path = "/account"
  async validate() {
    return this.makeJsonRequest<boolean>({
      method: "GET",
      url: `/validate`,
      headers: {
        "Content-Type": "application/json",
      },
      withCredentials: true,
    })
  }
}
