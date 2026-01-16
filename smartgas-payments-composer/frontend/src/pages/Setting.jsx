import {
  Box,
  Button,
  CircularProgress,
  Alert,
  FormControl,
  FormControlLabel,
  FormLabel,
  Grid,
  Radio,
  RadioGroup,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import {
  useGetSettingsQuery,
  useSetSettingMutation,
} from "../services/setting";
import { getSetting } from "../utils";
export default function SettingPage() {
  const { isLoading, data: settings, isError, error } = useGetSettingsQuery();

  const [
    addOrUpdateSetting,
    { isLoading: isLoadingCU, isError: isErrorCU, error: errorCU },
  ] = useSetSettingMutation();

  const [form, setForm] = useState({
    payment_provider: "",
    gas_pump_status: "",
    applicable_promotion_type: "none",
    ieps_regular: 0,
    ieps_premium: 0,
    ieps_diesel: 0,
  });

  const saveSetting = (settingName) => {
    addOrUpdateSetting({ name: settingName, value: form[settingName] });
  };

  useEffect(() => {
    if (settings?.length > 0) {
      setForm({
        payment_provider: getSetting(settings, "payment_provider"),
        gas_pump_status: getSetting(settings, "gas_pump_status"),
        applicable_promotion_type:
          getSetting(settings, "applicable_promotion_type") || "none",
        ieps_diesel: getSetting(settings, "ieps_diesel") || "0",
        ieps_premium: getSetting(settings, "ieps_premium") || "0",
        ieps_regular: getSetting(settings, "ieps_regular") || "0",
      });
    }
  }, [settings]);

  const handleChange = (e) => {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  if (isLoading || isLoadingCU) {
    return (
      <Box
        sx={{
          flex: 1,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
          rowGap: 3,
        }}
      >
        <CircularProgress />
        <Typography>Cargando...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ width: "100%" }}>
      {isErrorCU && <Alert severity="error">{errorCU}</Alert>}
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6}>
          <Box sx={{ flex: 1 }}>
            <FormControl sx={{ alignItems: "center", width: "100%" }}>
              <FormLabel id="demo-row-radio-buttons-group-label">
                Proveedor de pagos
              </FormLabel>
              <RadioGroup
                row
                aria-labelledby="demo-row-radio-buttons-group-label"
                name="payment_provider"
                value={form.payment_provider}
                onChange={handleChange}
              >
                <FormControlLabel
                  value="stripe"
                  control={<Radio />}
                  label="Stripe"
                />
                <FormControlLabel
                  value="swit"
                  control={<Radio />}
                  label="Swit"
                />
                <FormControlLabel
                  disabled
                  value=""
                  control={<Radio />}
                  label="Sin Valor"
                />
              </RadioGroup>
            </FormControl>
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("payment_provider")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
        <Grid item xs={12} sm={6}>
          <Box sx={{ flex: 1 }}>
            <FormControl sx={{ alignItems: "center", width: "100%" }}>
              <FormLabel id="demo-row-radio-buttons-group-label">
                Activacion de bombas
              </FormLabel>
              <RadioGroup
                row
                aria-labelledby="demo-row-radio-buttons-group-label"
                name="gas_pump_status"
                value={form.gas_pump_status}
                onChange={handleChange}
              >
                <FormControlLabel
                  value="enabled"
                  control={<Radio />}
                  label="Activadas"
                />
                <FormControlLabel
                  value="disabled"
                  control={<Radio />}
                  label="Desactivadas"
                />
                <FormControlLabel
                  disabled
                  value=""
                  control={<Radio />}
                  label="Sin Valor"
                />
              </RadioGroup>
            </FormControl>
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("gas_pump_status")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
      </Grid>

      <Grid container spacing={2} mt={3}>
        <Grid item xs={12} sm={4}>
          <Box
            sx={{
              width: "100%",
              textAlign: "center",
            }}
          >
            <TextField
              name="ieps_regular"
              label="IEPS Gasolina Regular"
              type="number"
              value={form.ieps_regular}
              onChange={handleChange}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("ieps_regular")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Box
            sx={{
              width: "100%",
              textAlign: "center",
            }}
          >
            <TextField
              name="ieps_premium"
              label="IEPS Gasolina Premium"
              type="number"
              value={form.ieps_premium}
              onChange={handleChange}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("ieps_premium")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Box
            sx={{
              width: "100%",
              textAlign: "center",
            }}
          >
            <TextField
              name="ieps_diesel"
              label="IEPS Diesel"
              type="number"
              onChange={handleChange}
              value={form.ieps_diesel}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("ieps_diesel")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
      </Grid>
      <Grid container spacing={2} mt={3}>
        <Grid item xs={12} sm={12}>
          <Box
            sx={{
              flex: 1,
              alignItems: "center",
              justifyContent: "center",
              display: "flex",
            }}
          >
            <FormControl sx={{ alignItems: "center", width: "100%" }}>
              <FormLabel id="demo-row-radio-buttons-group-label">
                Tipo de dinamica
              </FormLabel>
              <RadioGroup
                row
                aria-labelledby="demo-row-radio-buttons-group-label"
                name="applicable_promotion_type"
                value={form.applicable_promotion_type}
                onChange={handleChange}
              >
                <FormControlLabel
                  value="elegibility"
                  control={<Radio />}
                  label="Elegibilidad (Dinamica 2)"
                />
                <FormControlLabel
                  value="campaign"
                  control={<Radio />}
                  label="Campañas (Dinamica 1)"
                />
                <FormControlLabel
                  value="none"
                  control={<Radio />}
                  label="Ninguno"
                />
              </RadioGroup>
            </FormControl>
          </Box>
          <Box sx={{ textAlign: "center", marginTop: 2 }}>
            <Button
              variant="contained"
              onClick={() => saveSetting("applicable_promotion_type")}
            >
              Guardar
            </Button>
          </Box>
        </Grid>
      </Grid>
    </Box>
  );
}
