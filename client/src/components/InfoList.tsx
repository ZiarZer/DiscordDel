import styled from "styled-components";
import { Snowflake } from "./Snowflake";
import { useCallback } from "react";
import { Channel, Guild, InfoListFieldConfig, User } from "../types";

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

export function InfoList<T extends User | Guild | Channel>({
  currentObject,
  fields,
  getAvatarUrl,
}: {
  currentObject: Record<keyof T, string | number | undefined> | null;
  fields: Array<InfoListFieldConfig<T>>;
  getAvatarUrl: (param: T) => string | undefined;
}) {
  const defaultDisplay = useCallback((value?: string | number, isId = false) => {
    return isId ? <Snowflake snowflakeId={value} /> : 
    (value ?? '---')
  }, []);

  return (
    <Wrapper>
      <Avatar
        src={currentObject ? getAvatarUrl(currentObject as T) : undefined}
      />
      <ListWrapper>
        {fields.map(({ label, fieldName, display = defaultDisplay }, index) => {
          return (
            <span key={index}>
              <Bold>{label}:</Bold>&nbsp;
              {display(currentObject?.[fieldName], label.endsWith('ID'))}
            </span>
          );
        })}
      </ListWrapper>
    </Wrapper>
  );
}
