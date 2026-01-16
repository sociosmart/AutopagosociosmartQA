import {
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import React from "react";

interface TableProps {
  data: any[];
  columns: any[];
  total: number;
  paginationComponent: React.ReactNode;
}

export default function Table({
  data = [],
  columns = [],
  total,
  paginationComponent,
}: TableProps) {
  const [rows, _setRows] = React.useState(() => [...data]);

  const table = useReactTable({
    data: rows,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="w-full">
  <div className="w-full h-[65vh] md:h-[70vh] overflow-y-auto overflow-x-auto">
    <table className="table table-fixed min-w-full max-md:table-xs md:table-sm">
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <th key={header.id} className="px-6 py-3 md:w-[50px] w-[120px]">
                {header.isPlaceholder
                  ? null
                  : flexRender(header.column.columnDef.header, header.getContext())}
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody>
        {table.getRowModel().rows.map((row) => (
          <tr key={row.id}>
            {row.getVisibleCells().map((cell) => (
              <td key={cell.id} className="px-6 py-3 md:w-[50px] w-[120px] ">
                {flexRender(cell.column.columnDef.cell, cell.getContext())}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  </div>

  




      <div className="flex flex-col md:flex-row items-center justify-between mt-4 px-4 gap-2">

          <div className="md:w-full flex  gap-2  ">
            {paginationComponent && paginationComponent}
          </div>

          <div className="text-center  font-bold w-full md:w-auto">
            Total de Registros: {total}
          </div>

      </div>



    </div>
  );
  
}
