import { toast } from 'react-toastify';
import ApiService from "../Services/api.service";

export const getStats = () => () => {
  return ApiService.getStats().then(
    (data) => {
      console.log(data);
      if (data.error) {
          toast.error(data.error)
          return Promise.reject();
      }

      return Promise.resolve();
    },
    (error) => {
      const message =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString();
        console.log(message);

      toast.error(message);
      return Promise.reject();
    }
  );
};

