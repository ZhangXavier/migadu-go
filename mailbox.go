package migadu

import (
	"context"
	"net/http"
)

// Mailbox represents a mailbox in the Migadu API.
type Mailbox struct {
	Address               string     `json:"address,omitempty"`
	AutorespondActive     bool       `json:"autorespond_active,omitempty"`
	AutorespondBody       string     `json:"autorespond_body,omitempty"`
	AutorespondExpiresOn  string     `json:"autorespond_expires_on,omitempty"`
	AutorespondSubject    string     `json:"autorespond_subject,omitempty"`
	ChangedAt             string     `json:"changed_at,omitempty"`
	Delegations           []string   `json:"delegations,omitempty"`
	DomainName            string     `json:"domain_name,omitempty"`
	Expireable            bool       `json:"expireable,omitempty"`
	ExpiresOn             string     `json:"expires_on,omitempty"`
	FooterActive          bool       `json:"footer_active,omitempty"`
	FooterHTMLBody        string     `json:"footer_html_body,omitempty"`
	FooterPlainBody       string     `json:"footer_plain_body,omitempty"`
	Identities            []Identity `json:"identities,omitempty"`
	IsInternal            bool       `json:"is_internal,omitempty"`
	LastLoginAt           string     `json:"last_login_at,omitempty"`
	LocalPart             string     `json:"local_part,omitempty"`
	MayAccessImap         bool       `json:"may_access_imap,omitempty"`
	MayAccessManagesieve  bool       `json:"may_access_managesieve,omitempty"`
	MayAccessPop3         bool       `json:"may_access_pop3,omitempty"`
	MayReceive            bool       `json:"may_receive,omitempty"`
	MaySend               bool       `json:"may_send,omitempty"`
	Name                  string     `json:"name,omitempty"`
	Password              string     `json:"password,omitempty"`
	PasswordMethod        string     `json:"password_method,omitempty"`
	PasswordRecoveryEmail string     `json:"password_recovery_email,omitempty"`
	RecipientDenylist     []string   `json:"recipient_denylist,omitempty"`
	RemoveUponExpiry      bool       `json:"remove_upon_expiry,omitempty"`
	SenderAllowlist       []string   `json:"sender_allowlist,omitempty"`
	SenderDenylist        []string   `json:"sender_denylist,omitempty"`
	SpamAction            string     `json:"spam_action,omitempty"`
	SpamAggressiveness    string     `json:"spam_aggressiveness,omitempty"`
	StorageUsage          float64    `json:"storage_usage,omitempty"`
}

// ListMailboxes lists all the mailboxes for the domain configured on the client.
// Ir returns a pointer to an array of Mailbox structs and any error encountered.
func (c *Client) ListMailboxes(ctx context.Context) ([]*Mailbox, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddPath(MailboxesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest[struct {
		Mailboxes []*Mailbox `json:"mailboxes,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Mailboxes, nil
}

// GetMailbox retrieves a single mailbox given its local part name.
// It returns a pointer to a Mailbox struct and any error encountered.
func (c *Client) GetMailbox(ctx context.Context, localPart string) (*Mailbox, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddRestfulPath(MailboxesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Mailbox](c, ctx, req)
}

// NewMailbox creates a new mailbox given the local part, a display name, an invitation email and an optional password.
// An email will be sent to the email asking the user to set up a password. If a password is specified, the email will be used as the password recovery email.
// It returns a pointer to a Mailbox struct and any error encountered.
func (c *Client) NewMailbox(ctx context.Context, localPart string, displayName string, invitationEmail string, initialPassword string) (*Mailbox, error) {
	var mailbox = Mailbox{LocalPart: localPart, Name: displayName, PasswordRecoveryEmail: invitationEmail}
	if initialPassword != "" {
		mailbox.PasswordMethod = "password"
		mailbox.Password = initialPassword
	} else {
		mailbox.PasswordMethod = "invitation"
	}

	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPost).
		AddPath(MailboxesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(mailbox).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Mailbox](c, ctx, req)
}

// UpdateMailbox updates a mailbox in place given a pointer to a Mailbox struct.
// It returns a pointer to a new Mailbox struct and any error encountered.
func (c *Client) UpdateMailbox(ctx context.Context, localPart string, mb *Mailbox) (*Mailbox, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPut).
		AddRestfulPath(MailboxesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(mb).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Mailbox](c, ctx, req)
}

// DeleteMailbox deletes a mailbox given a pointer to a Mailbox struct.
// It returns any error encountered.
func (c *Client) DeleteMailbox(ctx context.Context, localPart string) error {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodDelete).
		AddRestfulPath(MailboxesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
