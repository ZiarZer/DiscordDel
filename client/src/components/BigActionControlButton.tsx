import styled from 'styled-components';
import { Action, ActionType, ActionScope } from '../types';
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
  action = null,
  running = false,
  onStopAction,
  onResumeLastAction,
}: {
  action?: Action | null;
  running?: boolean;
  onStopAction: () => void;
  onResumeLastAction: (
    type: ActionType,
    scope: ActionScope,
    targetId?: string
  ) => void;
}) {
  const backgroundColor = useMemo(() => {
    if (action == null) {
      return '#808080';
    }
    return running ? '#c23a22' : '#0f45d2';
  }, [action, running]);
  const { subtitle, mainText, onClick } = useMemo(() => {
    if (action == null) {
      return {
        subtitle: 'No action running',
        mainText: '_',
        onClick: () => {},
      };
    }
    const actionTitle =
      action?.targetId == null
        ? `${action.type} ${action.scope}`
        : `${action.type} ${action.scope} ${action.targetId}`;

    return {
      mainText: running ? 'Click to STOP action' : 'Click to RESUME action',
      subtitle: `${
        running ? 'Action running' : 'Last action run'
      }: ${actionTitle}`,
      onClick: running
        ? onStopAction
        : () => onResumeLastAction(action.type, action.scope, action.targetId),
    };
  }, [action, running]);

  return (
    <Wrapper $backgroundcolor={backgroundColor} onClick={onClick}>
      {subtitle}
      <SectionTitle>{mainText}</SectionTitle>
    </Wrapper>
  );
}
