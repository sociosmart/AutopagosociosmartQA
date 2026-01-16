import { format, parseISO } from "date-fns";

export const setAuthorizationHeaders = (headers, { getState }) => {
  const token = getState().auth.accessToken;

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  return headers;
};

export const formatDate = (date, fmtDate = "dd/MM/yyyy hh:mm aaaa") =>
  date ? format(date instanceof Date ? date : new Date(date), fmtDate) : "";

export const checkIfBetweenToDates = (from, to) =>
  new Date() >= new Date(from) && new Date() <= new Date(to);

export const gatherPurePermissions = (permissions) =>
  permissions.map((v) => v.name);

export const gatherPureGroups = (groups) => groups.map((v) => v.name);

export const gatherPermissions = (permissions, groups = []) => {
  let perms = gatherPurePermissions(permissions);

  groups.forEach(({ permissions }) => {
    let per = gatherPurePermissions(permissions);

    perms = [...perms, ...per];
  });

  return [...new Set(perms)];
};

export const checkPermissionInUser = (permission, user) =>
  user.is_admin ? true : !!user?.permissions.find((p) => permission === p);

export const getSetting = (settings, name) =>
  settings.find((s) => name === s.name)?.value || "";

export const getLastPaymentEvent = (events) =>
  !events.length ? "" : events[0].type;

export const dateToUTC = (date) => {
  let isoDate = date.toISOString();
  return `${isoDate.substring(0, 10)} ${isoDate.substring(11, 19)}`;
};

export const generateRandomNumber = (nPositions = 10) => {
  let ran = "";
  for (let i = 0; i < nPositions; i++) {
    ran += Math.floor(Math.random() * 10);
  }

  return ran;
};
