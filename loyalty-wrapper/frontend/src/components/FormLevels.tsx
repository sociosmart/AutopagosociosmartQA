import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";
import { BaseLevel } from "../models/base_level";

interface FormLevelsProps {
  level?: BaseLevel;
  onSubmit: (data: BaseLevel) => void;
}

const LevelSchema = Yup.object().shape({
  id: Yup.string().optional(),
  name: Yup.string()
    .required("El nombre es obligatorio")
    .max(100, "El nombre debe tener máximo 100 caracteres"),
  minPoints: Yup.number()
    .required("Los puntos mínimos son obligatorios")
    .min(0, "Debe ser mayor o igual a 0"),
  isActive: Yup.string()
    .oneOf(["true", "false"], "Estado inválido")
    .required("El estado es obligatorio"),
});

export default function FormLevels({ level, onSubmit }: FormLevelsProps) {
  const initialValues = {
    id: level?.id || "",
    name: level?.name || "",
    minPoints: level?.minPoints ?? 0,
    isActive: level ? (level.isActive ? "true" : "false") : "true",
  };

  return (
    <div className="flex justify-center items-center ">
      <Formik
        initialValues={initialValues}
        validationSchema={LevelSchema}
        enableReinitialize
        onSubmit={(values, { setSubmitting }) => {
          const levelData: BaseLevel = {
            id: values.id ? values.id : undefined,
            name: values.name,
            minPoints: values.minPoints,
            isActive: values.isActive === "true",
          };
          onSubmit(levelData);
          setSubmitting(false);
        }}
      >
        {({ isSubmitting }) => (
          <Form>
            <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
              <legend className="fieldset-legend">
                {level ? "Editar Nivel" : "Crear Nivel"}
              </legend>

              {/* Nombre */}
              <label className="fieldset-label" htmlFor="name">
                Nivel
              </label>
              <Field
                type="text"
                name="name"
                id="name"
                className="input"
                placeholder="Ingrese nombre del nivel"
              />
              <ErrorMessage
                name="name"
                component="div"
                className="text-red-500 mb-2"
              />

              {/* Puntos mínimos */}
              <label className="fieldset-label" htmlFor="minPoints">
                Puntos Mínimos
              </label>
              <Field
                type="number"
                name="minPoints"
                id="minPoints"
                className="input"
              />
              <ErrorMessage
                name="minPoints"
                component="div"
                className="text-red-500 mb-2"
              />

              {/* Estado */}
              <label className="fieldset-label" htmlFor="isActive">
                Estatus
              </label>
              <Field
                as="select"
                name="isActive"
                id="isActive"
                className="input"
              >
                <option value="true">Activo</option>
                <option value="false">Inactivo</option>
              </Field>
              <ErrorMessage
                name="isActive"
                component="div"
                className="text-red-500 mb-2"
              />

              <button
                type="submit"
                className="btn btn-primary mt-4"
                disabled={isSubmitting}
              >
                {level ? "Guardar Cambios" : "Crear Nivel"}
              </button>
            </fieldset>
          </Form>
        )}
      </Formik>
    </div>
  );
}
