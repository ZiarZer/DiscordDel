import styled from "styled-components";
import { Snowflake } from "./Snowflake";

const ListWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-self: start;
  padding: 0 1em;
  font-size: 0.9em;
  text-align: left;
`;
const Bold = styled.span`
  font-weight: 700;
`;

const Avatar = styled.img`
  height: 4em;
  width: 4em;
  border-radius: 4px;
  opacity: ${({ src }) => (src ? 1 : 0)};
  transition: opacity 150ms ease-in-out;
`;

const Wrapper = styled.div`
  display: flex;
  align-items: center;
`;

export function InfoList({
  currentObject,
  fields,
  getAvatarUrl,
}: {
  data: Array<{ label: string; value: string }>;
}) {
  return (
    <Wrapper>
      <Avatar src={currentObject ? getAvatarUrl(currentObject) : undefined} />
      <ListWrapper>
        {fields.map(({ label, fieldName }) => {
          return (
            <span key={fieldName}>
              <Bold>{label}:</Bold>&nbsp;
              {label.endsWith("ID") ? (
                <Snowflake snowflakeId={currentObject?.[fieldName]} />
              ) : (
                currentObject?.[fieldName] ?? "---"
              )}
            </span>
          );
        })}
      </ListWrapper>
    </Wrapper>
  );
}
