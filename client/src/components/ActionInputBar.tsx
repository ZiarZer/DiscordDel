import { ChangeEvent, useRef, useState } from 'react';

import styled from 'styled-components';
import { Button as BaseButton, Input as BaseInput } from '.';

const Button = styled(BaseButton)`
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
`;

const Input = styled(BaseInput)`
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
`;

const Label = styled.label`
  align-self: center;
`;

const Wrapper = styled.div`
  display: flex;
  align-items: stretch;
`;

export function ActionInputBar({
  inputPlaceholder,
  buttonText,
  enabled = false,
  secret = false,
  onEdit,
  onSubmit,
}: {
  inputPlaceholder: string
  buttonText: string
  enabled?: boolean
  secret?: boolean
  onSubmit: (param: string) => void
  onEdit: (e: ChangeEvent) => void
}) {
  const [showSecret, setShowSecret] = useState(!secret);
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <Wrapper>
      <Input
        type={showSecret ? 'text' : 'password'}
        placeholder={inputPlaceholder}
        ref={inputRef}
        onChange={onEdit}
      />
      <Button disabled={!enabled} onClick={() => onSubmit(inputRef.current?.value ?? '')}>
        {buttonText}
      </Button>
      {secret ? (
        <>
          <input
            type='checkbox'
            id='display-token-checkbox'
            onChange={() => setShowSecret(!showSecret)}
          />
          <Label htmlFor='display-token-checkbox'>Display token</Label>
        </>
      ) : null}
    </Wrapper>
  );
}
