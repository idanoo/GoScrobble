import ApiService from "../Services/api.service";

export const getStats = () => {
  return ApiService.getStats().then(
    (data) => {
      return data.data;
    }
  );
};

