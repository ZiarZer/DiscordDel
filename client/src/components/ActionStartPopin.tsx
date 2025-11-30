import { useCallback, useMemo, useState } from 'react';
import { Popin } from './Popin';
import { Button } from '.';
import styled from 'styled-components';
import { ActionScope, ActionType } from '../types/actions';

const ActionOptions = {
  CRAWL: {
    crawlReactions: 'Crawl reactions',
  },
  DELETE: {
    deletePinned: 'Delete pinned messages',
    deleteThreadFirstMessage: 'Delete thread first message',
  },
};

const ButtonsWrapper = styled.div`
  display: flex;
  gap: 1em;
`;
const Checkbox = styled.input.attrs({ type: 'checkbox' })``;

type ActionStartPopinProps = {
  type: ActionType;
  scope: ActionScope;
  targetId?: string;
  onStartAction: (
    type: ActionType,
    scope: ActionScope,
    targetId?: string,
    options?: object
  ) => void;
  onClose: () => void;
  isOpen: boolean;
};

export function ActionStartPopin({
  type,
  scope,
  targetId,
  onStartAction,
  isOpen,
  onClose,
}: ActionStartPopinProps) {
  const actionText = useMemo(
    () => (type === 'CRAWL' ? 'Crawl' : 'Delete crawled data of'),
    [type]
  );
  const targetText = useMemo(() => {
    if (scope === 'CHANNEL') {
      return `channel ${targetId}`;
    }
    if (scope === 'GUILD') {
      return `guild ${targetId}`;
    }
    return 'all guilds';
  }, [scope, targetId]);

  const [requestOption, setRequestOption] = useState({});

  const handleStartAction = useCallback(() => {
    onStartAction(type, scope, targetId, requestOption);
    onClose();
  }, [type, scope, targetId, requestOption, onClose, onStartAction]);

  return (
    <Popin
      onClose={onClose}
      isOpen={isOpen}
      title={`${actionText} ${targetText}`}
    >
      {Object.entries(ActionOptions[type]).map(([key, label]) => (
        <label>
          <Checkbox
            onChange={(e) =>
              setRequestOption({ ...requestOption, [key]: e.target.checked })
            }
          />
          {label}
        </label>
      ))}
      <ButtonsWrapper>
        <Button onClick={handleStartAction}>Start action</Button>
        <Button onClick={onClose}>Cancel</Button>
      </ButtonsWrapper>
    </Popin>
  );
}
