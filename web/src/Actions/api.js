import ApiService from "../Services/api.service";

export const getStats = () => {
  return ApiService.getStats().then(
    (data) => {
      return data.data;
    }
  );
};

export const getRecentScrobbles = (id) => {
  return ApiService.getRecentScrobbles(id).then(
    (data) => {
      return data.data;
    }
  );
};

