import { BaseProduct } from "../models/base_product";
import { BaseGasStation } from "../models/base_gas_station";
import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";
import { BaseMargin } from "../models/base_margin";

interface FormMarginsProps {
  margin?: BaseMargin;
  products: BaseProduct[];
  gasStations: BaseGasStation[];
  onSubmit: (data: BaseMargin) => void;
}



const MarginSchema = Yup.object().shape({
  id: Yup.string().optional(),
  margin: Yup.number()
    .required("El margen es obligatorio")
    .min(0, "El margen debe ser un número positivo"),
  points: Yup.number()
    .required("Los puntos son obligatorios")
    .min(0, "Los puntos deben ser un número positivo"),
  product: Yup.object().shape({
    id: Yup.string().required("Debe seleccionar un producto"),
  }),
});

export default function FormMargins({ margin, products, gasStations, onSubmit }: FormMarginsProps) {
  const initialValues = {
    id: margin?.id || "",
    marginType: margin?.marginType , 
    margin: margin?.margin || 0,
    points: margin?.points || 0,
    product: margin?.product || { id: "", name: "", isActive: false },
    gasStation: margin?.gasStation ||{ id: "", name: "", externalId: "", crePermission: "", latitude: "", longitude: "" }
  };
  return (
    <div className="flex justify-center items-center">
      <Formik
        initialValues={initialValues}
        validationSchema={MarginSchema}
        enableReinitialize
        onSubmit={(values, { setSubmitting }) => {
            const marginData: BaseMargin = {
              id: values.id || "",
              marginType: values.marginType || "by_margin", 
              margin: values.margin,
              points: values.points,
              product: values.product.id as unknown as BaseProduct, 
              gasStation: values.gasStation.id ? values.gasStation.id as unknown as BaseGasStation : null
            };

          
            onSubmit(marginData);
            setSubmitting(false);
          }}
          
          
      >
        {({ isSubmitting, values }) => (
          <Form>
            <fieldset className="fieldset w-md bg-base-200 border border-base-300 p-14 rounded-box">
              <legend className="fieldset-legend">
                {margin ? "Editar Margen" : "Crear Margen"}
              </legend>





              <label className="fieldset-label" htmlFor="marginType">
                Tipo
              </label>
              <Field as="select" name="marginType" id="marginType" className="input">
                <option value="by_margin">Margen</option>
                <option value="by_liter">Litro</option>
              </Field>
              <ErrorMessage name="marginType" component="div" className="text-red-500 mb-2" />


              {values.marginType !== "by_liter" && (
                 <>
                <label className="fieldset-label" htmlFor="margin">
                  Margen (%)
                </label>
                <Field
                  type="number"
                  name="margin"
                  id="margin"
                  className="input"
                  placeholder="Ingrese el margen"
                />
                <ErrorMessage name="margin" component="div" className="text-red-500 mb-2" />

                </>
          )}


                <label className="fieldset-label" htmlFor="points">
                  Puntos
                </label>
                <Field
                  type="number"
                  name="points"
                  id="points"
                  className="input"
                  placeholder="Ingrese los puntos"
                />
                <ErrorMessage name="points" component="div" className="text-red-500 mb-2" />

                 




              <label className="fieldset-label" htmlFor="product">
                Producto
              </label>
              <Field as="select" name="product.id" id="product" className="input">
                <option value="">Seleccione un producto</option>
                {products.map((product) => (
                  <option key={product.id} value={product.id}>
                    {product.name}
                  </option>
                ))}
              </Field>

              <ErrorMessage name="product.id" component="div" className="text-red-500 mb-2" />


              <label className="fieldset-label" htmlFor="gasStation">
                Estación de Gas
              </label>
              <Field as="select" name="gasStation.id" id="gasStation" className="input">
                <option value="">NINGUNO</option>
                {gasStations.map((gasStation) => (
                  <option key={gasStation.id} value={gasStation.id}>
                    {gasStation.name}
                  </option>
                ))}
              </Field>

              <ErrorMessage name="gasStation.id" component="div" className="text-red-500 mb-2" />

              <button type="submit" className="btn btn-primary mt-4" disabled={isSubmitting}>
                {margin ? "Guardar Cambios" : "Crear Margen"}
              </button>
            </fieldset>
          </Form>
        )}
      </Formik>
    </div>
  );
}
