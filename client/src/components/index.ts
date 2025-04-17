import styled from "styled-components";

const Button = styled.button`
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.25em 0.5em;
  font-size: 0.9em;
  font-weight: 500;
  font-family: inherit;
  background-color: #1a1a1a;
  cursor: pointer;
  transition: border-color 0.25s;
  color: inherit;
  &:hover {
    border-color: #646cff;
  }
  &:hover:disabled {
    border-color: transparent;
  }
  &:disabled {
    cursor: default;
    background-color: #282828;
    color: #777777;
  }
`;

const Input = styled.input`
  border-radius: 8px;
  font-size: 0.9em;
  padding: 0.6em 1.2em;
  width: 14em;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
`;

export {
  Button,
  Input,
};
