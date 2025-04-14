import styled from "styled-components";
import { Snowflake } from "./Snowflake";

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-self: start;
  padding: 0 1em;
  font-size: 0.9em;
`;
const Bold = styled.span`
  font-weight: 700;
`;

export function InfoList({
  data = [],
}: {
  data: Array<{ label: string; value: string }>;
}) {
  return (
    <Wrapper>
      {data.map(({ label, value }) => {
        const displayedValue = value ?? "---";
        return (
          <span key={label}>
            <Bold>{label}:</Bold>&nbsp;
            {label.endsWith('ID') ? (
              <Snowflake snowflakeId={displayedValue}></Snowflake>
            ) : (
              displayedValue
            )}
          </span>
        );
      })}
    </Wrapper>
  );
}
