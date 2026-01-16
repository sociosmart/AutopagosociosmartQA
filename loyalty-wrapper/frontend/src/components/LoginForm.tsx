import { useMutation } from "@apollo/client";
import { Formik, Field, Form, ErrorMessage } from "formik";
import * as Yup from "yup";
import { POST_LOGIN } from "../screens/Login";
import { LoginError } from "../models/login_error";
import { LoginSuccess } from "../models/login_success";
import { useDispatch } from "react-redux";
import { logIn } from "../store/authSlice";
import { useNavigate } from "react-router-dom";

const validationSchema = Yup.object({
  email: Yup.string()
    .email("Correo en formato incorrecto")
    .required("Required"),
  password: Yup.string()
    .min(6, "La contraseña necesita al menos 6 caracteres")
    .required("Required"),
});

export default function FormLogin() {
  const navigate = useNavigate();

  const [login, { loading, error }] = useMutation<{
    login: LoginSuccess | LoginError;
  }>(POST_LOGIN);
  const dispatch = useDispatch();

  const handleSubmit = async (values: { email: string; password: string }) => {
    try {
      const { data } = await login({
        variables: {
          email: values.email,
          password: values.password,
        },
      });

      if (data?.login && "__typename" in data?.login) {
        if (data?.login.__typename === "LoginSuccess") {
          const LoginSuccess = data?.login as LoginSuccess;

          dispatch(
            logIn({
              access_token: LoginSuccess.accessToken,
              refresh_token: LoginSuccess.refreshToken,
            }),
          );

          navigate("/protected/users");
        } else if (data?.login.__typename === "LoginError") {
          const LoginError = data?.login as LoginError;
          if (LoginError.type === "invalid_credentials") {
          }
        }
      }
    } catch (err) {
      // TODO: Implement an alert here
    }
  };
  return (
    <div className="flex justify-center items-center min-h-screen">
      <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
        <legend className="fieldset-legend">Login</legend>

        <Formik
          initialValues={{ email: "", password: "" }}
          validationSchema={validationSchema}
          onSubmit={handleSubmit}
        >
          {({ isSubmitting }) => (
            <Form>
              <div>
                <label className="fieldset-label" htmlFor="email">
                  Email
                </label>
                <Field name="email" className="input" placeholder="Email" />
                <ErrorMessage
                  name="email"
                  component="div"
                  className="text-red-500 text-sm"
                />
              </div>

              <div>
                <label className="fieldset-label" htmlFor="password">
                  Password
                </label>
                <Field
                  type="password"
                  name="password"
                  className="input"
                  placeholder="Password"
                />
                <ErrorMessage
                  name="password"
                  component="div"
                  className="text-red-500 text-sm"
                />
              </div>

              <button
                type="submit"
                className="btn btn-neutral mt-4"
                disabled={isSubmitting || loading}
              >
                {loading ? "Logging in..." : "Login"}
              </button>

              {error && (
                <div className="text-red-500 mt-2">
                  {error.message || "An error occurred"}
                </div>
              )}
            </Form>
          )}
        </Formik>
      </fieldset>
    </div>
  );
}
