import { useCallback } from "react";
import styled from "styled-components";

const Code = styled.code<{ $snowflakeId: string | number | null }>`
  cursor: ${({ $snowflakeId }) => $snowflakeId ? 'copy' : 'default'};
`;

export function Snowflake({ snowflakeId = null }: { snowflakeId?: string | number | null }) {
  const copyId = useCallback(
    () => snowflakeId == null ? () => {} : navigator.clipboard.writeText(snowflakeId.toString()),
    [snowflakeId]
  );
  return (
    <Code onClick={copyId} $snowflakeId={snowflakeId}>
      #{snowflakeId ?? '---'}
    </Code>
  );
}
