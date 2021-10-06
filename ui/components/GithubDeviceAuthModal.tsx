import { CircularProgress } from "@material-ui/core";
import * as React from "react";
import styled from "styled-components";
import useAuth from "../hooks/auth";
import { GetGithubDeviceCodeResponse } from "../lib/api/applications/applications.pb";
import { GitProviderName } from "../lib/types";
import Alert from "./Alert";
import Button from "./Button";
import Flex from "./Flex";
import Icon, { IconType } from "./Icon";
import Modal from "./Modal";
import Text from "./Text";

type Props = {
  className?: string;
  open: boolean;
  onSuccess: (token: string) => void;
  onClose: () => void;
  repoName: string;
};

const Pad = styled(Flex)`
  padding: 8px 0;
`;

const ModalContent = styled(({ codeRes, onSuccess, onError, className }) => {
  // Move this to a component so that we get the cancel logic when the modal closes.
  const { getGithubAuthStatus } = useAuth();
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    if (!codeRes) {
      return;
    }
    setLoading(true);
    const { cancel, promise } = getGithubAuthStatus(codeRes);

    promise
      .then((authRes) => {
        onSuccess(authRes.accessToken);
      })
      .catch((err) => {
        onError(err);
      })
      .finally(() => setLoading(false));

    return cancel;
  }, [codeRes]);

  return (
    <div className={className}>
      <Pad wide center>
        <Text size="extraLarge">{codeRes.userCode}</Text>
      </Pad>
      <Pad wide center>
        <a target="_blank" href={codeRes.validationURI}>
          <Button type="button" variant="contained" color="primary">
            <Flex align>
              Authorize Github Access{" "}
              <Icon size="base" type={IconType.ExternalTab} />
            </Flex>
          </Button>
        </a>
      </Pad>
      <Pad wide center>
        {loading && <div>Waiting for authorization to be completed...</div>}
      </Pad>
    </div>
  );
})`
  ${Icon} {
    margin-left: 8px;
  }
`;

function GithubDeviceAuthModal({
  className,
  open,
  onClose,
  repoName,
  onSuccess,
}: Props) {
  const [codeRes, setCodeRes] =
    React.useState<GetGithubDeviceCodeResponse>(null);
  const { getGithubDeviceCode, storeProviderToken } = useAuth();
  const [codeLoading, setCodeLoading] = React.useState(true);
  const [error, setError] = React.useState(null);

  React.useEffect(() => {
    if (!open) {
      return;
    }

    setCodeLoading(true);

    getGithubDeviceCode()
      .then((res) => {
        setCodeRes(res);
      })
      .finally(() => setCodeLoading(false));
  }, [open]);
  return (
    <Modal
      className={className}
      title="Authenticate with Github"
      open={open}
      onClose={onClose}
      description={`Weave GitOps needs to authenitcate with the Git Provider for the ${repoName} repo`}
    >
      <p>
        Paste this code into the Github Device Activation field to grant Weave
        GitOps temporary access:
      </p>
      {error && (
        <Alert severity="error" title="Error" message={error.message} />
      )}

      <Flex wide center>
        {codeLoading || !codeRes ? (
          <CircularProgress />
        ) : (
          <ModalContent
            onSuccess={(token) => {
              storeProviderToken(GitProviderName.GitHub, token);
              onSuccess(token);
              onClose();
            }}
            onError={(err) => setError(err)}
            codeRes={codeRes}
          />
        )}
      </Flex>
    </Modal>
  );
}

export default styled(GithubDeviceAuthModal)``;
