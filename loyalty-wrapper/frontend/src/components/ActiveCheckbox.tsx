import { FaTimes } from "react-icons/fa";
import { FaCheck } from "react-icons/fa";

interface ActiveCheckboxProps {
  isActive: boolean;
}

export default function ActiveCheckbox({ isActive }: ActiveCheckboxProps) {
  return (
    <div>
      {isActive ? (
        <FaCheck className="text-500" />
      ) : (
        <FaTimes className="text-500" />
      )}
    </div>
  );
}
