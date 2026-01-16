import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";
import { BaseBenefit } from "../models/base_benefit";
import { BaseLevel } from "../models/base_level";

interface FormBenefitsProps {
  benefit?: BaseBenefit;
  levels: BaseLevel[];
  onSubmit: (data: Omit<BaseBenefit, "createdAt" | "updatedAt">) => void;
}

const BenefitSchema = Yup.object().shape({
  id: Yup.string().optional(),
  name: Yup.string()
    .required("El nombre es obligatorio")
    .max(100, "El nombre debe tener máximo 100 caracteres"),
  type: Yup.string().required("El tipo es obligatorio"),
  frequency: Yup.string().when("type", (type, schema) => {
    if (typeof type === "string" && type !== "gas" && type !== "periferics") {
      return schema.required("La frecuencia es obligatoria");
    }
    return schema.notRequired();
  }),
  discount: Yup.number()
    .required("El descuento es obligatorio")
    .min(0, "Debe ser mayor o igual a 0"),
  stock: Yup.number().when("type", (type, schema) => {
    if (typeof type === "string" && type !== "gas" && type !== "periferics") {
      return schema.required("El stock es obligatorio").min(0, "Debe ser mayor o igual a 0");
    }
    return schema.notRequired();
  }),
  numTimes: Yup.number().when("frequency", (frequency, schema) => {
    if (typeof frequency === "string" && frequency === "n_times") {
      return schema.required("El número de veces es obligatorio").min(0, "Debe ser mayor o igual a 0");
    }
    return schema.notRequired();
  }),
  isActive: Yup.string()
    .oneOf(["true", "false"], "Estado inválido")
    .required("El estado es obligatorio"),
  level: Yup.object().shape({
    id: Yup.string().required("El nivel es obligatorio"),
  }),
  externalProductId: Yup.string().when("type", (type, schema) => {
    if (typeof type === "string" && type !== "gas" && type !== "periferics") {
      return schema.required("El ID del producto es obligatorio");
    }
    return schema.notRequired();
  }),
});



export default function FormBenefits({ benefit, levels, onSubmit }: FormBenefitsProps) {
  const initialValues = {
    id: benefit?.id || "",
    name: benefit?.name || "",
    type: benefit?.type || "digital",
    frequency: benefit?.frequency || "always",
    discount: benefit?.discount ?? 0,
    stock: benefit?.stock ?? 0,
    externalProductId: benefit?.externalProductId || "",
    numTimes: benefit?.numTimes ?? 0,
    isActive: benefit ? (benefit.isActive ? "true" : "false") : "true",
    level: benefit?.level || { id: "", name: "", minPoints: 0, isActive: true },
    dependency: benefit ? (benefit.dependency ? "true" : "false") : "true",
    minAmount: benefit?.minAmount ?? 0,
  };

  return (
    
    <div className="flex justify-center items-center">
      <Formik
        initialValues={initialValues}
        validationSchema={BenefitSchema}
        enableReinitialize
        onSubmit={(values, { setSubmitting }) => {
          const benefitData: Omit<BaseBenefit, "createdAt" | "updatedAt"> = {
            id: values.id ? values.id : "",
            name: values.name,
            type: values.type,
            frequency: values.frequency,
            discount: values.discount,
            stock: values.stock,
            externalProductId: values.externalProductId,
            numTimes: values.numTimes,
            isActive: values.isActive === "true",
            level: { id: values.level.id, name: values.level.name } as BaseLevel,
            dependency: values.dependency  === "true",
            minAmount: values.minAmount
          };
          onSubmit(benefitData);
          setSubmitting(false);
        }}
      >
        {({ isSubmitting, values }) => (
          <Form>
            <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
              <legend className="fieldset-legend">
                {benefit ? "Editar Beneficio" : "Crear Beneficio"}
              </legend>


              {/* Tipo - select */}
              <label className="fieldset-label" htmlFor="type">
                Tipo
              </label>
              <Field as="select" name="type" id="type" className="input">
                <option value="digital">Digital</option>
                <option value="physical">Físico</option>
                <option value="gas">Combustible</option>
                <option value="periferics">Periféricos</option>
              </Field>
              <ErrorMessage name="type" component="div" className="text-red-500 mb-2" />



              {/* Nombre */}
              <label className="fieldset-label" htmlFor="name">
                Nombre
              </label>
              <Field
                type="text"
                name="name"
                id="name"
                className="input"
                placeholder="Ingrese el nombre del beneficio"
              />
              <ErrorMessage name="name" component="div" className="text-red-500 mb-2" />

              {/* Producto Externo */}
              {values.type !== "gas" && values.type !== "periferics" && (
                <>
                  <label className="fieldset-label" htmlFor="externalProductId">
                    Id del Producto Externo
                  </label>
                  <Field
                    type="text"
                    name="externalProductId"
                    id="externalProductId"
                    className="input"
                    placeholder="Ingrese el ID del Producto"
                  />
                  <ErrorMessage name="externalProductId" component="div" className="text-red-500 mb-2" />
                </>
              )}

              {/* Frecuencia - select */}
              {values.type !== "gas" && values.type !== "periferics" && (
                <>
                  <label className="fieldset-label" htmlFor="frequency">
                    Frecuencia
                  </label>
                  <Field as="select" name="frequency" id="frequency" className="input">
                    <option value="always">Siempre</option>
                    <option value="n_times">N veces</option>
                    <option value="hourly">Cada Hora</option>
                    <option value="daily">Diariamente</option>
                    <option value="weekly">Semanalmente</option>
                    <option value="monthly">Mensualmente</option>
                  </Field>
                  <ErrorMessage name="frequency" component="div" className="text-red-500 mb-2" />
                </>
              )}

              {/* Descuento */}
              {(values.type === "gas" || values.type === "periferics") && (
                <>
                  <label className="fieldset-label" htmlFor="discount">
                    Descuento
                  </label>
                  <Field type="number" name="discount" id="discount" className="input" />
                  <ErrorMessage name="discount" component="div" className="text-red-500 mb-2" />
                </>
              )}


              {/* Stock */}
              {values.type !== "gas" && values.type !== "periferics" && (
                <>
                  <label className="fieldset-label" htmlFor="stock">
                    Stock
                  </label>
                  <Field type="number" name="stock" id="stock" className="input" />
                  <ErrorMessage name="stock" component="div" className="text-red-500 mb-2" />
                </>
              )}

              {/* Número de veces */}
              {values.frequency === "n_times" && (
                <>
                  <label className="fieldset-label" htmlFor="numTimes">
                    Número de Veces
                  </label>
                  <Field type="number" name="numTimes" id="numTimes" className="input" />
                  <ErrorMessage name="numTimes" component="div" className="text-red-500 mb-2" />
                </>
              )}

                {/* MIN AMOUNT */}
                  {values.type === "physical" && (
                <>
                  <label className="fieldset-label" htmlFor="minAmount">
                    Cantidad Minima
                  </label>
                  <Field type="number" name="minAmount" id="minAmount" className="input" />
                  <ErrorMessage name="minAmount" component="div" className="text-red-500 mb-2" />
                </>
              )}

              {/* Nivel */}
              <label className="fieldset-label" htmlFor="level">
                Nivel
              </label>
              <Field as="select" name="level.id" id="level" className="input">
                <option value="">Seleccione un nivel</option>
                {levels.map((level) => (
                  <option key={level.id} value={level.id}>
                    {level.name}
                  </option>
                ))}
              </Field>
              <ErrorMessage name="level.id" component="div" className="text-red-500 mb-2" />

              {/* Dependency */}
              {values.type === "physical" && (
                <>
              <label className="fieldset-label" htmlFor="dependency">
                Es Dependiente
              </label>
              <Field as="select" name="dependency" id="dependency" className="input">
                <option value="true">Si</option>
                <option value="false">No</option>
              </Field>
              <ErrorMessage name="dependency" component="div" className="text-red-500 mb-2" />
              </>
              )}

              {/* Estado */}
              <label className="fieldset-label" htmlFor="isActive">
                Estatus
              </label>
              <Field as="select" name="isActive" id="isActive" className="input">
                <option value="true">Activo</option>
                <option value="false">Inactivo</option>
              </Field>
              <ErrorMessage name="isActive" component="div" className="text-red-500 mb-2" />

              <button type="submit" className="btn btn-primary mt-4" disabled={isSubmitting}>
                {benefit ? "Guardar Cambios" : "Crear Beneficio"}
              </button>
            </fieldset>
          </Form>
        )}
      </Formik>
    </div>
  );
}
