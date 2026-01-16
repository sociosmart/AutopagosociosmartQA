import { Box } from "@mui/system";
import {
  useAddCampaignMutation,
  useEditCampaignMutation,
  useGetCampaignQuery,
} from "../services/campaigns";
import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useMemo } from "react";
import {
  Alert,
  CircularProgress,
  Grid,
  Typography,
  TextField,
  FormControlLabel,
  Checkbox,
  Autocomplete,
} from "@mui/material";
import * as Yup from "yup";
import { LoadingButton } from "@mui/lab";
import { useFormik } from "formik";
import { useGetGasStationsAllQuery } from "../services/gasStations";
import { MobileDateTimePicker } from "@mui/x-date-pickers/MobileDateTimePicker";
import { dateToUTC } from "../utils";
import { useDispatch } from "react-redux";
import { showSnackbar } from "../store/snackbar";

const validationSchema = Yup.object({
  name: Yup.string()
    .min(2, "Nombre demasiado corto")
    .required("Nombres es requerido"),
  discount: Yup.number().positive("Tienes que introducir un numero positivo"),
  valid_from: Yup.date().required("Requerido"),
  valid_to: Yup.date()
    .min(
      Yup.ref("valid_from"),
      "Tienes que elegir una fecha mayor al inicio de vigencia",
    )
    .required("Requerido"),
});

export default function CampaignDetail() {
  const { id } = useParams();
  //
  const isEditing = useMemo(() => Boolean(id), [id]);

  const navigate = useNavigate();
  const dispatch = useDispatch();

  const {
    data: campaignDetail,
    isLoading: isLoadingCampaignDetail,
    isError: isErrorLoadingCampaign,
    error: campaignError,
  } = useGetCampaignQuery(id, {
    skip: !isEditing,
  });

  const [
    editCampaign,
    {
      isLoading: isLoadingEditing,
      isError: isEditingCampaignError,
      error: errorEditingCampaign,
      status: editCampaignStatus,
    },
  ] = useEditCampaignMutation();

  const [
    addCampaign,
    {
      isLoading: isLoadingAdding,
      isError: isAddingCampaignError,
      error: errorAddingCampaign,
      status: addCampaignStatus,
    },
  ] = useAddCampaignMutation();

  const { data: gas_stations = [], isLoading: isLoadingGasStations } =
    useGetGasStationsAllQuery();

  const formik = useFormik({
    initialValues: {
      name: "",
      discount: 0,
      valid_from: new Date(),
      valid_to: new Date(),
      active: true,
      gas_stations: [],
    },
    validationSchema,
    onSubmit: (values) => {
      let parsedValues = {
        ...values,
        valid_from: dateToUTC(values["valid_from"]),
        valid_to: dateToUTC(values["valid_to"]),
      };

      if (!isEditing) {
        addCampaign(parsedValues);
      } else {
        editCampaign({ id, ...parsedValues });
      }
    },
  });

  useEffect(() => {
    if (addCampaignStatus === "fulfilled") {
      navigate("/protected/campaigns");
      dispatch(
        showSnackbar({
          message: "Campaña agregada",
        }),
      );
    }
  }, [addCampaignStatus]);

  useEffect(() => {
    if (editCampaignStatus === "fulfilled") {
      navigate("/protected/campaigns");
      dispatch(
        showSnackbar({
          message: "Campaña editada",
        }),
      );
    }
  }, [editCampaignStatus]);

  useEffect(() => {
    if (campaignDetail) {
      Object.keys(campaignDetail).forEach((key) => {
        if (key === "valid_from" || key === "valid_to")
          formik.setFieldValue(key, new Date(campaignDetail[key]));
        else formik.setFieldValue(key, campaignDetail[key]);
      });
    }
  }, [campaignDetail]);

  if (isErrorLoadingCampaign) {
    return (
      <Box>
        <Alert severity="error">{campaignError}</Alert>
      </Box>
    );
  }

  if (isLoadingCampaignDetail || isLoadingGasStations) {
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
      {isAddingCampaignError && (
        <Box my={3}>
          <Alert severity="error">{errorAddingCampaign}</Alert>
        </Box>
      )}
      {isEditingCampaignError && (
        <Box my={3}>
          <Alert severity="error">{errorEditingCampaign}</Alert>
        </Box>
      )}
      <Box component="form" onSubmit={formik.handleSubmit}>
        <Grid container spacing={2}>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              name="name"
              variant="outlined"
              label="Nombre de la campaña"
              required
              value={formik.values.name}
              onChange={formik.handleChange}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
            />
          </Grid>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              name="discount"
              variant="outlined"
              label="Descuento por litro"
              type="number"
              required
              value={formik.values.discount}
              onChange={formik.handleChange}
              error={formik.touched.discount && Boolean(formik.errors.discount)}
              helperText={formik.touched.discount && formik.errors.discount}
            />
          </Grid>
        </Grid>
        <Box my={2} />
        <Grid container spacing={2}>
          <Grid item md={6} xs={12}>
            <MobileDateTimePicker
              label="Valido desde"
              value={formik.values.valid_from}
              onChange={(value) => {
                formik.setFieldValue("valid_from", value);
              }}
              slotProps={{
                textField: {
                  required: true,
                  error:
                    formik.touched.valid_from &&
                    Boolean(formik.errors.valid_from),
                  helperText:
                    formik.touched.valid_from && formik.errors.valid_from,
                },
              }}
            />
          </Grid>
          <Grid item md={6} xs={12}>
            <MobileDateTimePicker
              label="Valido hasta"
              value={formik.values.valid_to}
              onChange={(value) => {
                formik.setFieldValue("valid_to", value);
              }}
              slotProps={{
                textField: {
                  required: true,
                  error:
                    formik.touched.valid_to && Boolean(formik.errors.valid_to),
                  helperText: formik.touched.valid_to && formik.errors.valid_to,
                },
              }}
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
            name="active"
            value={formik.values.active}
            onChange={formik.handleChange}
            control={<Checkbox checked={formik.values.active} />}
            label="Activo"
          />
        </Box>
        <Grid container spacing={2}>
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
        <Box sx={{ textAlign: "center", my: 3 }}>
          <LoadingButton
            loading={isLoadingAdding || isLoadingEditing}
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
