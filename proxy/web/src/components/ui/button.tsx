import type React from "react";

type ButtonProps = React.ComponentPropsWithRef<"button">;

export function Button({ children }: ButtonProps) {
  return <button>{children}</button>;
}
