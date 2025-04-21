import { ReactNode } from "react";

export function StatusMessage({
  message,
  success = false,
}: {
  message: string | ReactNode;
  success: boolean;
}) {
  return (
    <div>
      {success ? (
        <span style={{ color: "green" }}>&#x2714;</span>
      ) : (
        <span>&#x274C;</span>
      )}
      &nbsp;
      {message}
    </div>
  );
}
