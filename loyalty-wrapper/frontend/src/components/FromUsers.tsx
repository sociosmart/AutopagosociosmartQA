import React from "react";
import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";
import { BaseUser } from "../models/base_user";


interface FormUserProps {
  user?: BaseUser; 
  onSubmit: (data: BaseUser) => void;
}

const UserSchema = Yup.object().shape({
    id: Yup.string().optional(),
    firstName: Yup.string()
      .required("El nombre es obligatorio")
      .max(100, "El nombre debe tener máximo 100 caracteres"),
    lastName: Yup.string()
      .required("Los apellidos son obligatorios")
      .max(100, "Los apellidos deben tener máximo 100 caracteres"),
    email: Yup.string()
      .email("Correo en formato incorrecto")
      .required("El correo es obligatorio"),
    password: Yup.string().when("id", {
      is: (id: string | undefined) => !id, 
      then: (schema) =>
        schema
          .required("La contraseña es obligatoria")
          .min(6, "La contraseña debe tener al menos 6 caracteres")
          .max(20, "La contraseña no puede tener más de 20 caracteres"),
      otherwise: (schema) =>
        schema
          .notRequired()
          .test(
            "password-length",
            "La contraseña debe tener entre 6 y 20 caracteres",
            (value) => !value || (value.length >= 6 && value.length <= 20)
          ),
    }),
    isActive: Yup.string()
      .oneOf(["true", "false"], "Estado inválido")
      .required("El estado es obligatorio"),
  });
  


export default function FormUsers({ user, onSubmit }: FormUserProps) {
  const initialValues = {
    id: user?.id || "",
    firstName: user?.firstName || "",
    lastName: user?.lastName || "",
    email: user?.email || "",
    password: "",
    isActive: user ? (user.isActive ? "true" : "false") : "true",
  };

  return (
    <div className="flex justify-center items-center">
      <Formik
        initialValues={initialValues}
        validationSchema={UserSchema}
        enableReinitialize={true}
        onSubmit={(values, { setSubmitting }) => {
          const userData: BaseUser = {
            id: values.id ? values.id : undefined,
            firstName: values.firstName,
            lastName: values.lastName,
            email: values.email,
            password: values.password,
            isActive: values.isActive === "true",
          };
          onSubmit(userData);
          setSubmitting(false);
        }}
      >
        {({ isSubmitting }) => (
          <Form>
            <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
              <legend className="fieldset-legend">
                {user ? "Editar Usuario" : "Crear Usuario"}
              </legend>

              <label className="fieldset-label" htmlFor="firstName">
                Nombres
              </label>
              <Field
                type="text"
                name="firstName"
                id="firstName"
                className="input"
                placeholder="Nombre(s)"
              />
              <ErrorMessage name="name" component="div" className="text-red-500 mb-2" />

              <label className="fieldset-label" htmlFor="lastName">
                Apellidos
              </label>
              <Field
                type="text"
                name="lastName"
                id="lastName"
                className="input"
                placeholder="Apellidos"
              />
              <ErrorMessage name="name" component="div" className="text-red-500 mb-2" />


                <label className="fieldset-label" htmlFor="email">
                  Email
                </label>
                <Field      
                type="email"
                name="email"
                id="email"
                className="input"
                placeholder="Email"   
                />
                <ErrorMessage name="email" component="div"  className="text-red-500 mb-2" />


                <label className="fieldset-label" htmlFor="password">
                    Password
                </label>
                <Field
                    type="password"
                    name="password"
                    id="password"
                    className="input"
                />
                <ErrorMessage name="password" component="div" className="text-red-500 mb-2" />





              <label className="fieldset-label" htmlFor="isActive">
                Estatus
              </label>
              <Field as="select" name="isActive" id="isActive" className="input">
                <option value="true">Activo</option>
                <option value="false">Inactivo</option>
              </Field>
              <ErrorMessage name="isActive" component="div" className="text-red-500 mb-2" />

              <button type="submit" className="btn btn-primary mt-4" disabled={isSubmitting}>
                {user ? "Guardar Cambios" : "Crear Usuario"}
              </button>
              
            </fieldset>
          </Form>
        )}
      </Formik>
    </div>
  );
}
