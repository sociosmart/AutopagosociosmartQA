import { Alert, Box, Container, Paper, TextField, Typography } from "@mui/material"
import LoadingButton from '@mui/lab/LoadingButton';
import { useFormik } from "formik"
import * as Yup from "yup"

import SmartGasLogo from "../assets/smartgas_logo.png"
import { useLogInMutation } from "../services/auth"

const validationSchema = Yup.object({
  email: Yup.string()
    .email('Correo invalido')
    .required('Correo requerido'),
  password: Yup.string()
    .min(3, "Contraseña demasiado corta")
    .max(100, "Contraseña demasiado larga")
})

export default function () {

  const [
    logIn,
    {
      isLoading,
      isError,
      error
    }
  ] = useLogInMutation()


  const formik = useFormik({
    initialValues: {
      email: "",
      password: ""
    },
    validationSchema,
    onSubmit: values => {
      logIn(values)
    }
  })


  return (
    <Container maxWidth="xl">
      <Box height={"100vh"} sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: 'center',
        alignItems: "center",
        backgroundColor: "theme.secondary"
      }}>
        <Box maxWidth="xs">
          <img src={SmartGasLogo} alt="SmartGas Log" width="100%" />
        </Box>
        <Paper elevation={14}>
          <Box m={5}>
            <Typography variant="h5" component="h6" textAlign="center">
              Log in
            </Typography>

            {
              isError &&
              <Box mt={3}><Alert severity="error">{error}</Alert></Box>
            }
            <Box component="form" mt={3} onSubmit={formik.handleSubmit}>
              <TextField
                id="email"
                name="email"
                variant="outlined"
                label="Correo"
                required
                fullWidth
                value={formik.values.email}
                onChange={formik.handleChange}
                error={formik.touched.email && Boolean(formik.errors.email)}
                helperText={formik.touched.email && formik.errors.email}
              />
              <TextField
                name="password"
                type="password"
                variant="outlined"
                label="Contraseña"
                required
                fullWidth
                value={formik.values.password}
                onChange={formik.handleChange}
                error={formik.touched.password && Boolean(formik.errors.password)}
                helperText={formik.touched.password && formik.errors.password}
                sx={{
                  mt: 3
                }}
              />
              <LoadingButton
                loading={isLoading}
                type="submit"
                fullWidth
                variant="contained"
                sx={{
                  mt: 3
                }}
              >
                Log in
              </LoadingButton>
            </Box>
          </Box>
        </Paper>
      </Box>
    </Container>
  )
}
