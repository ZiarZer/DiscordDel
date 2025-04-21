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
  flex: 1;
`;

const Label = styled.label`
  align-self: center;
`;

const Wrapper = styled.div`
  display: flex;
  align-items: stretch;
`;

type ActionInputBarProps = {
  inputPlaceholder: string;
  buttonText: string;
  enabled?: boolean;
  secret?: boolean;
  onSubmit: () => void;
  onChange: (e: ChangeEvent) => void;
};

export function ActionInputBar({
  inputPlaceholder,
  buttonText,
  enabled = false,
  secret = false,
  onSubmit,
  onChange = () => {},
}: ActionInputBarProps) {
  const [showSecret, setShowSecret] = useState(!secret);
  const inputRef = useRef<HTMLInputElement>(null);
  const id = `display-checkbox-${Math.floor(Math.random() * 10)}`;

  return (
    <Wrapper>
      <Input
        type={showSecret ? 'text' : 'password'}
        placeholder={inputPlaceholder}
        ref={inputRef}
        onChange={onChange}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onSubmit();
          }
        }}
      />
      <Button disabled={!enabled} onClick={onSubmit}>
        {buttonText}
      </Button>
      {secret ? (
        <>
          <input
            type='checkbox'
            id={id}
            onChange={() => setShowSecret(!showSecret)}
            style={{ marginLeft: '1em' }}
          />
          <Label htmlFor={id}>Display</Label>
        </>
      ) : null}
    </Wrapper>
  );
}
