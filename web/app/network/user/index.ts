import login from "./post/login";
import refresh from "./post/refresh";

const userServices = {
  get: {

  },
  post: {
    login,
    refresh
  }
}

export default userServices;