import { Button, Stack } from "@mui/material";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { StoragePolicy } from "../../../../../api/dashboard";
import { PolicyType } from "../../../../../api/explorer";
import { DenseFilledTextField } from "../../../../Common/StyledComponents";
import SettingForm from "../../../../Pages/Setting/SettingForm";
import { NoMarginHelperText } from "../../../Settings/Settings";
import { AddWizardProps } from "../../AddWizardDialog";

const Cloud189Wizard = ({ onSubmit }: AddWizardProps) => {
  const { t } = useTranslation("dashboard");
  const formRef = useRef<HTMLFormElement>(null);
  const [policy, setPolicy] = useState<StoragePolicy>({
    id: 0,
    name: "",
    type: PolicyType.cloud189,
    dir_name_rule: "uploads/{uid}/{path}",
    file_name_rule: "{uuid}_{originname}",
    settings: {
      chunk_size: 10 << 20, // 10MB
    },
    edges: {},
  });

  const handleSubmit = () => {
    if (!formRef.current?.checkValidity()) {
      formRef.current?.reportValidity();
      return;
    }
    onSubmit(policy);
  };

  return (
    <form ref={formRef} onSubmit={handleSubmit}>
      <Stack spacing={2}>
        <SettingForm title={t("policy.name")} lgWidth={12}>
          <DenseFilledTextField
            fullWidth
            required
            value={policy.name}
            onChange={(e) => setPolicy({ ...policy, name: e.target.value })}
          />
          <NoMarginHelperText>{t("policy.policyName")}</NoMarginHelperText>
        </SettingForm>
        <SettingForm title={t("policy.accessCredential")} lgWidth={12}>
          <DenseFilledTextField
            placeholder={t("policy.cloud189Account")}
            fullWidth
            required
            value={policy.access_key || ""}
            onChange={(e) => setPolicy({ ...policy, access_key: e.target.value })}
          />
          <NoMarginHelperText>{t("policy.cloud189AccountDes")}</NoMarginHelperText>
          <DenseFilledTextField
            placeholder={t("policy.cloud189Password")}
            type="password"
            sx={{ mt: 1 }}
            fullWidth
            required
            value={policy.secret_key || ""}
            onChange={(e) => setPolicy({ ...policy, secret_key: e.target.value })}
          />
          <NoMarginHelperText>{t("policy.cloud189PasswordDes")}</NoMarginHelperText>
        </SettingForm>
      </Stack>
      <Button variant="contained" color="primary" sx={{ mt: 2 }} onClick={handleSubmit}>
        {t("policy.create")}
      </Button>
    </form>
  );
};

export default Cloud189Wizard;
