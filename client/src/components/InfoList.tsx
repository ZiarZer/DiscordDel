import styled from "styled-components";
import { Snowflake } from "./Snowflake";
import { ReactElement, useCallback } from "react";

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

type FieldConfig = {
  fieldName: string;
  label: string;
  display: (value?: string | number, isId?: boolean) => string | ReactElement
}

export function InfoList({
  currentObject,
  fields,
  getAvatarUrl,
}: {
  currentObject: Record<string, string | number>;
  fields: Array<FieldConfig>;
  getAvatarUrl: (currentObject: object) => string;
}) {
  const defaultDisplay = useCallback((value?: string | number, isId = false) => {
    return isId ? <Snowflake snowflakeId={value} /> : 
    (value ?? '---')
  }, []);

  return (
    <Wrapper>
      <Avatar src={currentObject ? getAvatarUrl(currentObject) : undefined} />
      <ListWrapper>
        {fields.map(
          ({
            label,
            fieldName,
            display = defaultDisplay,
          }) => {
            return (
              <span key={fieldName}>
                <Bold>{label}:</Bold>&nbsp;
                {display(currentObject?.[fieldName], label.endsWith('ID'))}
              </span>
            );
          }
        )}
      </ListWrapper>
    </Wrapper>
  );
}
