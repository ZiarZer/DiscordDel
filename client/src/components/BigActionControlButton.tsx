import styled from 'styled-components';
import { Log } from './Log';
import { Action } from '../types';
import { useMemo } from 'react';

const Wrapper = styled.div<{ $backgroundcolor: string }>`
  background-color: ${({ $backgroundcolor }) => $backgroundcolor};
  color: white;
  border-radius: 1em;
  padding: 1em;
  text-align: center;
  display: flex;
  flex-direction: column;
  width: 100%;
  font-family: monospace;
  transition: transform ease-in-out 150ms, background-color ease-in-out 150ms;
  &:active {
    transform: scale(0.95);
  }
  cursor: pointer;
`;

const SectionTitle = styled.h2`
  margin: 0;
`;

export function BigActionControlButton({
  action,
  onClick,
}: {
  action?: Action;
  onClick: () => void;
}) {
  const actionTitle = useMemo(() => {
    if (action == null) {
      return null;
    }
    let result = `${action.type} ${action.scope}`;
    if (action.targetId != null) {
      result += ` ${action.targetId}`;
    }
    return result;
  }, [action]);
  const backgroundColor = useMemo(
    () => (action == null ? '#808080' : '#c23a22'),
    [action]
  );
  const subtitle = useMemo(
    () =>
      action == null ? 'No action running' : `Action running: ${actionTitle}`,
    [action]
  );
  const mainText = useMemo(
    () => (action == null ? '_' : 'Click to STOP action'),
    [action]
  );

  return (
    <Wrapper $backgroundcolor={backgroundColor} onClick={onClick}>
      {subtitle}
      <SectionTitle>{mainText}</SectionTitle>
    </Wrapper>
  );
}
