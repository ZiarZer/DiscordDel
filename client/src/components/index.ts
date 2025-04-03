import styled from "styled-components";

const Button = styled.button`
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.6em 1.2em;
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
  }
`;

export {
  Button
}
