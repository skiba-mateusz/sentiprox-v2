import type React from "react";

type ButtonVariant = "primary" | "secondary" | "circle";
type ButtonSize = "small" | "medium" | "large";

interface ButtonProps {
  children: React.ReactNode;
  variant?: ButtonVariant;
  size?: ButtonSize;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
  disabled?: boolean;
}

export function Button({
  variant = "primary",
  size = "medium",
  onClick,
  disabled,
  children,
}: ButtonProps) {
  return (
    <button
      className={`btn btn--${variant} btn--${size}`}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
}
