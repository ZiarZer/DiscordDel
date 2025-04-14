import { useCallback } from "react";
import styled from "styled-components";

const Code = styled.code<{ $snowflakeId: string|null }>`
  cursor: ${({ $snowflakeId }) => $snowflakeId ? 'copy' : 'default'};
`;

export function Snowflake({ snowflakeId = null }: { snowflakeId?: string|null }) {
  const copyId = useCallback(
    () => snowflakeId == null ? () => {} : navigator.clipboard.writeText(snowflakeId),
    [snowflakeId]
  );
  return (
    <Code onClick={copyId} $snowflakeId={snowflakeId}>
      #{snowflakeId ?? '---'}
    </Code>
  );
}
