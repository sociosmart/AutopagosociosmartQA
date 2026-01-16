import React from "react";
import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";

export interface Product {
  id?: string;
  name: string;
  codename: string;
  isActive: boolean;
}

interface FormProductsProps {
  product?: Product; 
  onSubmit: (data: Product) => void;
}

const ProductSchema = Yup.object().shape({
  id: Yup.string().optional(),
  name: Yup.string()
    .required("El nombre es obligatorio")
    .max(100, "El nombre debe tener máximo 100 caracteres"),
    codename: Yup.string()
    .required("El codigo es obligatorio"),
  isActive: Yup.string()
    .oneOf(["true", "false"], "Estado inválido")
    .required("El estado es obligatorio"),
});


export default function FormProducts({ product, onSubmit }: FormProductsProps) {
  const initialValues = {
    id: product?.id || "",
    name: product?.name || "",
    codename: product?.codename || "",
    isActive: product ? (product.isActive ? "true" : "false") : "true",
  };

  return (
    <div className="flex justify-center items-center ">
      <Formik
        initialValues={initialValues}
        validationSchema={ProductSchema}
        enableReinitialize={true}
        onSubmit={(values, { setSubmitting }) => {
          const productData: Product = {
            id: values.id ? values.id : undefined,
            name: values.name,
            codename: values.codename,
            isActive: values.isActive === "true",
          };
          onSubmit(productData);
          setSubmitting(false);
        }}
      >
        {({ isSubmitting }) => (
          <Form>
            <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
              <legend className="fieldset-legend">
                {product ? "Editar Producto" : "Crear Producto"}
              </legend>

              <label className="fieldset-label" htmlFor="name">
                Nombre del Producto
              </label>
              <Field
                type="text"
                name="name"
                id="name"
                className="input"
                placeholder="Ingrese nombre del producto"
              />
              <ErrorMessage name="name" component="div" className="text-red-500 mb-2" />

              <label className="fieldset-label" htmlFor="codename">
                Codigo del Producto
              </label>
              <Field
                type="text"
                name="codename"
                id="codename"
                className="input"
                placeholder="Ingrese codigo del producto"
              />
              <ErrorMessage name="codename" component="div" className="text-red-500 mb-2" />

              <label className="fieldset-label" htmlFor="isActive">
                Estatus
              </label>
              <Field as="select" name="isActive" id="isActive" className="input">
                <option value="true">Activo</option>
                <option value="false">Inactivo</option>
              </Field>
              <ErrorMessage name="isActive" component="div" className="text-red-500 mb-2" />

              <button type="submit" className="btn btn-primary mt-4" disabled={isSubmitting}>
                {product ? "Guardar Cambios" : "Crear Producto"}
              </button>
              
            </fieldset>
          </Form>
        )}
      </Formik>
    </div>
  );
}
