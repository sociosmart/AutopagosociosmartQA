import { useFormik } from "formik";
import * as Yup from "yup";
import { useNavigate, useParams } from "react-router-dom";
import { Box } from "@mui/system";
import {
  Alert,
  Autocomplete,
  Checkbox,
  CircularProgress,
  FormControlLabel,
  Grid,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useMemo } from "react";
import {
  useGetAllGroupsQuery,
  useGetAllPermissionsQuery,
} from "../services/permissions";
import {
  useAddUserMutation,
  useEditUserMutation,
  useGetUserDetailQuery,
} from "../services/user";
import { LoadingButton } from "@mui/lab";
import { useDispatch } from "react-redux";
import { showSnackbar } from "../store/snackbar";
import { useGetGasStationsAllQuery } from "../services/gasStations";

const validationSchema = Yup.object({
  email: Yup.string().email("Correo invalido").required("Correo requerido"),
  first_name: Yup.string()
    .min(2, "Nombre demasiado corto")
    .required("Nombres es requerido"),
  last_name: Yup.string()
    .min(2, "Apellidos demasiado corto")
    .required("Apellidos requerido"),
  password: Yup.string()
    .min(3, "Contraseña demasiado corta")
    .max(100, "Contraseña demasiado larga"),
});

const permissionsToHuman = {
  view_gas_stations: "Ver Estaciones",
  edit_gas_station: "Editar Estaciones",
  add_gas_station: "Agregar estacions",
  can_do_payment_actions: "Permiso para reembolso y preset",
  view_gas_pumps: "Ver Bombas",
  add_gas_pump: "Agregar bombas",
  edit_gas_pump: "Editar bombas",
  view_users: "Ver usuarios",
  edit_user: "Editar usuarios",
  view_payments: "Ver pagos y cargas",
  view_synchronizations: "Ver sincronizaciones",
  add_synchronizations: "Poder sincronizar bombas",
  view_campaigns: "Ver Campañas",
  edit_campaign: "Editar Campañas",
  add_campaign: "Agregar Campañas",
  view_elegibility_levels: "Ver niveles",
  edit_elegibility_level: "Editar niveles",
  add_elegibility_level: "Agregar niveles",
  view_customer_levels: "Ver niveles de clientes",
  edit_customer_level: "Editar manualmente nivel de usuarios",
  add_customer_level: "Agregar manualmente nivel de usuarios",
  view_all_customers: "Ver todos los clientes",
  view_all_elegibility_levels: "Ver todos los niveles",
};

export default function UserDetail() {
  const { id } = useParams();
  //
  const isEditing = useMemo(() => Boolean(id), [id]);

  const navigate = useNavigate();

  const dispatch = useDispatch();

  const { data: permissions = [], isLoading: isLoadingPermissions } =
    useGetAllPermissionsQuery();

  const { data: groups = [], isLoading: isLoadingGroups } =
    useGetAllGroupsQuery();

  const { data: gas_stations = [], isLoading: isLoadingGasStations } =
    useGetGasStationsAllQuery();

  const [
    addUser,
    {
      isLoading: isAdding,
      isError: isErrorAdding,
      error: errorAdding,
      status: addUserStatus,
    },
  ] = useAddUserMutation();

  const [
    editUser,
    {
      isLoading: isEditting,
      isError: isEdittingError,
      error: errorEditting,
      status: editUserStatus,
    },
  ] = useEditUserMutation();

  useEffect(() => {
    if (addUserStatus === "fulfilled") {
      navigate("/protected/users");
      dispatch(
        showSnackbar({
          message: "Usuario agregado",
        }),
      );
    }
  }, [addUserStatus]);

  useEffect(() => {
    if (editUserStatus === "fulfilled") {
      navigate("/protected/users");
      dispatch(
        showSnackbar({
          message: "Usuario Editado",
        }),
      );
    }
  }, [editUserStatus]);

  const {
    data: user,
    isLoading: loadingUser,
    isError: isErrorLoadingUser,
    error: errorLoadingUser,
  } = useGetUserDetailQuery(id, { skip: !isEditing });

  const formik = useFormik({
    initialValues: {
      email: "",
      password: "",
      first_name: "",
      last_name: "",
      is_admin: false,
      active: true,
      groups: [],
      permissions: [],
      gas_stations: [],
    },
    validationSchema,
    onSubmit: (values) => {
      if (!values.password) {
        delete values["password"];
      }
      if (!isEditing) {
        addUser(values);
      } else {
        editUser({ id, ...values });
      }
    },
  });

  useEffect(() => {
    if (user) {
      Object.keys(user).forEach((key) => {
        formik.setFieldValue(key, user[key]);
      });
    }
  }, [user]);

  if (isErrorLoadingUser) {
    return (
      <Box>
        <Alert severity="error">{errorLoadingUser}</Alert>
      </Box>
    );
  }

  if (
    isLoadingGroups ||
    isLoadingPermissions ||
    loadingUser ||
    isLoadingGasStations
  ) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
          flex: 1,
        }}
      >
        <CircularProgress />
        <Typography variant="body2" align="center" mt={3}>
          Cargando...
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {isErrorAdding && (
        <Box my={3}>
          <Alert severity="error">{errorAdding}</Alert>
        </Box>
      )}
      {isEdittingError && (
        <Box my={3}>
          <Alert severity="error">{errorEditting}</Alert>
        </Box>
      )}
      <Box component="form" onSubmit={formik.handleSubmit}>
        <Grid container spacing={2}>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              name="first_name"
              variant="outlined"
              label="Nombres"
              required
              value={formik.values.first_name}
              onChange={formik.handleChange}
              error={
                formik.touched.first_name && Boolean(formik.errors.first_name)
              }
              helperText={formik.touched.first_name && formik.errors.first_name}
            />
          </Grid>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              name="last_name"
              variant="outlined"
              label="Apellidos"
              required
              value={formik.values.last_name}
              onChange={formik.handleChange}
              error={
                formik.touched.last_name && Boolean(formik.errors.last_name)
              }
              helperText={formik.touched.last_name && formik.errors.last_name}
            />
          </Grid>
        </Grid>
        <Box sx={{ my: 2 }} />
        <Grid container spacing={2}>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              name="email"
              variant="outlined"
              label="Correo"
              required={!isEditing}
              disabled={isEditing}
              value={formik.values.email}
              onChange={formik.handleChange}
              error={formik.touched.email && Boolean(formik.errors.email)}
              helperText={formik.touched.email && formik.errors.email}
            />
          </Grid>
          <Grid item md={6} xs={12}>
            <TextField
              name="password"
              type="password"
              variant="outlined"
              label="Contraseña"
              required={!isEditing}
              fullWidth
              value={formik.values.password}
              onChange={formik.handleChange}
              error={formik.touched.password && Boolean(formik.errors.password)}
              helperText={
                formik.touched.password && formik.errors.password
                  ? formik.errors.password
                  : isEditing
                  ? "Si escribes una Contraseña, la anterior sera remplazada"
                  : ""
              }
            />
          </Grid>
        </Grid>
        <Box
          sx={{
            display: "flex",
            flexDirection: "row",
            justifyContent: "space-around",
            my: 2,
          }}
        >
          <FormControlLabel
            name="is_admin"
            value={formik.values.is_admin}
            onChange={formik.handleChange}
            control={<Checkbox checked={formik.values.is_admin} />}
            label="Admin"
          />
          <FormControlLabel
            name="active"
            value={formik.values.active}
            onChange={formik.handleChange}
            control={<Checkbox checked={formik.values.active} />}
            label="Activo"
          />
        </Box>
        {!formik.values.is_admin && (
          <Grid container spacing={2}>
            <Grid item md={6} xs={12}>
              <Autocomplete
                multiple
                options={groups || []}
                onChange={(_, value) => formik.setFieldValue("groups", value)}
                isOptionEqualToValue={(option, value) => option.id === value.id}
                value={formik.values.groups}
                getOptionLabel={(option) => option.name}
                defaultValue={[]}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    variant="outlined"
                    label="Grupos"
                    placeholder="Buscar grupo..."
                  />
                )}
              />
            </Grid>
            <Box sx={{ my: 3 }} />
            <Grid item md={6} xs={12}>
              <Autocomplete
                multiple
                onChange={(_, value) =>
                  formik.setFieldValue("permissions", value)
                }
                isOptionEqualToValue={(option, value) => option.id === value.id}
                value={formik.values.permissions}
                options={permissions || []}
                getOptionLabel={(option) =>
                  permissionsToHuman[option.name] || option.name
                }
                defaultValue={[]}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    variant="outlined"
                    label="Permisos"
                    placeholder="Permisos individuales"
                  />
                )}
              />
            </Grid>
            <Grid item md={6} xs={12}>
              <Autocomplete
                multiple
                options={gas_stations || []}
                onChange={(_, value) =>
                  formik.setFieldValue("gas_stations", value)
                }
                isOptionEqualToValue={(option, value) => option.id === value.id}
                value={formik.values.gas_stations}
                getOptionLabel={(option) => option.name}
                defaultValue={[]}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    variant="outlined"
                    label="Estaciones"
                    placeholder="Buscar Estaciones"
                  />
                )}
              />
            </Grid>
          </Grid>
        )}
        <Box sx={{ textAlign: "center", my: 3 }}>
          <LoadingButton
            loading={isAdding | isEditting}
            type="submit"
            variant="contained"
          >
            Guardar
          </LoadingButton>
        </Box>
      </Box>
    </Box>
  );
}
