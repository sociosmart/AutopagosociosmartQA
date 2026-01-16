interface PaginationProps {
  onNextPage: () => void;
  onPreviousPage: () => void;
  hasNext: boolean;
  hasPrevious: boolean;
}

export default function Pagination({
  onNextPage,
  onPreviousPage,
  hasNext,
  hasPrevious,
}: PaginationProps) {
  return (
    <div className="flex justify-center md:justify-start gap-2 mt-5 px-4 md:px-0 w-auto">
      <button
        className="btn btn-outline btn-md w-full sm:w-1/2 md:w-auto"
        onClick={onPreviousPage}
        disabled={!hasPrevious}
      >
        Anterior
      </button>
      <button
        className="btn btn-outline btn-md w-full sm:w-1/2 md:w-auto"
        onClick={onNextPage}
        disabled={!hasNext}
      >
        Siguiente
      </button>
    </div>
  );
}


