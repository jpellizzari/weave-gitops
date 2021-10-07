package capi

import (
	"fmt"
	"io"
	"strings"
)

// TemplatesRetriever defines the interface that adapters
// need to implement in order to return an array of templates.
type TemplatesRetriever interface {
	Source() string
	RetrieveTemplates() ([]Template, error)
	RetrieveTemplatesByProvider(provider string) ([]Template, error)
	RetrieveTemplateParameters(name string) ([]TemplateParameter, error)
}

// TemplateRenderer defines the interface that adapters
// need to implement in order to render a template populated
// with parameter values.
type TemplateRenderer interface {
	RenderTemplateWithParameters(name string, parameters map[string]string, creds Credentials) (string, error)
}

// TemplatePullRequester defines the interface that adapters
// need to implement in order to create a pull request from
// a CAPI template. Implementers should return the web URI of
// the pull request.
type TemplatePullRequester interface {
	CreatePullRequestFromTemplate(params CreatePullRequestFromTemplateParams) (string, error)
}

// CredentialsRetriever defines the interface that adapters
// need to implement in order to retrieve CAPI credentials.
type CredentialsRetriever interface {
	Source() string
	RetrieveCredentials() ([]Credentials, error)
}

type Template struct {
	Name        string
	Description string
	Provider    string
}

type TemplateParameter struct {
	Name        string
	Description string
	Required    bool
	Options     []string
}

type Credentials struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string
}

type CreatePullRequestFromTemplateParams struct {
	TemplateName    string
	ParameterValues map[string]string
	RepositoryURL   string
	HeadBranch      string
	BaseBranch      string
	Title           string
	Description     string
	CommitMessage   string
	Credentials     Credentials
}

// GetTemplates uses a TemplatesRetriever adapter to show
// a list of templates to the console.
func GetTemplates(r TemplatesRetriever, w io.Writer) error {
	ts, err := r.RetrieveTemplates()
	if err != nil {
		return fmt.Errorf("unable to retrieve templates from %q: %w", r.Source(), err)
	}

	if len(ts) > 0 {
		fmt.Fprintf(w, "NAME\tPROVIDER\tDESCRIPTION\n")

		for _, t := range ts {
			fmt.Fprintf(w, "%s", t.Name)
			fmt.Fprintf(w, "\t%s", t.Provider)

			if t.Description != "" {
				fmt.Fprintf(w, "\t%s", t.Description)
			}

			fmt.Fprintln(w, "")
		}

		return nil
	}

	fmt.Fprintf(w, "No templates found.\n")

	return nil
}

// GetTemplatesByProvider uses a TemplatesRetriever adapter to show
// a list of templates for a given provider to the console.
func GetTemplatesByProvider(provider string, r TemplatesRetriever, w io.Writer) error {
	ts, err := r.RetrieveTemplatesByProvider(provider)
	if err != nil {
		return fmt.Errorf("unable to retrieve templates from %q: %w", r.Source(), err)
	}

	if len(ts) > 0 {
		fmt.Fprintf(w, "NAME\tPROVIDER\tDESCRIPTION\n")

		for _, t := range ts {
			fmt.Fprintf(w, "%s", t.Name)
			fmt.Fprintf(w, "\t%s", t.Provider)

			if t.Description != "" {
				fmt.Fprintf(w, "\t%s", t.Description)
			}

			fmt.Fprintln(w, "")
		}

		return nil
	}

	fmt.Fprintf(w, "No templates were found for provider %q.\n", provider)

	return nil
}

// GetTemplateParameters uses a TemplatesRetriever adapter
// to show a list of parameters for a given template.
func GetTemplateParameters(name string, r TemplatesRetriever, w io.Writer) error {
	ps, err := r.RetrieveTemplateParameters(name)
	if err != nil {
		return fmt.Errorf("unable to retrieve parameters for template %q from %q: %w", name, r.Source(), err)
	}

	if len(ps) > 0 {
		fmt.Fprintf(w, "NAME\tREQUIRED\tDESCRIPTION\tOPTIONS\n")

		for _, t := range ps {
			fmt.Fprintf(w, "%s", t.Name)
			fmt.Fprintf(w, "\t%t", t.Required)

			if t.Description != "" {
				fmt.Fprintf(w, "\t%s", t.Description)
			}

			if t.Options != nil {
				optionsStr := strings.Join(t.Options, ", ")
				fmt.Fprintf(w, "\t%s", optionsStr)
			}

			fmt.Fprintln(w, "")
		}

		return nil
	}

	fmt.Fprintf(w, "No template parameters were found.")

	return nil
}

// RenderTemplate uses a TemplateRenderer adapter to show
// a template populated with parameter values.
func RenderTemplateWithParameters(name string, parameters map[string]string, creds Credentials, r TemplateRenderer, w io.Writer) error {
	t, err := r.RenderTemplateWithParameters(name, parameters, creds)
	if err != nil {
		return fmt.Errorf("unable to render template %q: %w", name, err)
	}

	if t != "" {
		fmt.Fprint(w, t)
		return nil
	}

	fmt.Fprintf(w, "No template found.")

	return nil
}

func CreatePullRequestFromTemplate(params CreatePullRequestFromTemplateParams, r TemplatePullRequester, w io.Writer) error {
	res, err := r.CreatePullRequestFromTemplate(params)
	if err != nil {
		return fmt.Errorf("unable to create pull request: %w", err)
	}

	fmt.Fprintf(w, "Created pull request: %s\n", res)

	return nil
}

// GetCredentials uses a CredentialsRetriever adapter to show
// a list of CAPI credentials.
func GetCredentials(r CredentialsRetriever, w io.Writer) error {
	cs, err := r.RetrieveCredentials()
	if err != nil {
		return fmt.Errorf("unable to retrieve credentials from %q: %w", r.Source(), err)
	}

	if len(cs) > 0 {
		fmt.Fprintf(w, "NAME\tINFRASTRUCTURE PROVIDER\n")

		for _, c := range cs {
			fmt.Fprintf(w, "%s", c.Name)
			// Extract the infra provider name from ClusterKind
			provider := c.Kind[:strings.Index(c.Kind, "Cluster")]
			fmt.Fprintf(w, "\t%s", provider)
			fmt.Fprintln(w, "")
		}

		return nil
	}

	fmt.Fprintf(w, "No credentials found.")

	return nil
}
