package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/template"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	activityResourceSuccess = "success"
	activityResourceFailure = "failure"
)

// ActivityResourceLoaderFunc retrieves the specified resource string value.
type ActivityResourceLoaderFunc func(context.Context, string) (string, error)

// ActivityEntityReference refers to an entity involved in the Activity
type ActivityEntityReference struct {
	Id   string
	Name string
	Type string
}

// ActivityBuilder automates the construction of Activity Feed messsages.
type ActivityBuilder struct {
	Description struct {
		Text           string
		Template       string
		Resource       string
		ResourceLoader ActivityResourceLoaderFunc
	}
	ServiceType string
	Action      string
	Severity    struct {
		Display string
		Event   string
	}
	Operands struct {
		Object ActivityEntityReference
		Target ActivityEntityReference
		Tenant ActivityEntityReference
		User   ActivityEntityReference
	}
}

func (b *ActivityBuilder) WithDescriptionText(text string) *ActivityBuilder {
	b.Description.Text = text
	return b
}

func (b *ActivityBuilder) WithDescriptionResource(resource string) *ActivityBuilder {
	b.Description.Resource = resource
	return b
}

func (b *ActivityBuilder) WithDescriptionResourceLoader(loader ActivityResourceLoaderFunc) *ActivityBuilder {
	b.Description.ResourceLoader = loader
	return b
}

func (b *ActivityBuilder) WithDescriptionTemplate(template string) *ActivityBuilder {
	b.Description.Template = template
	return b
}

func (b *ActivityBuilder) WithObject(object ActivityEntityReference) *ActivityBuilder {
	b.Operands.Object = object
	return b
}

func (b *ActivityBuilder) WithTenant(tenant ActivityEntityReference) *ActivityBuilder {
	b.Operands.Tenant = tenant
	return b
}

func (b *ActivityBuilder) WithTarget(target ActivityEntityReference) *ActivityBuilder {
	b.Operands.Target = target
	return b
}

func (b *ActivityBuilder) WithServiceType(serviceType string) *ActivityBuilder {
	b.ServiceType = serviceType
	return b
}

func (b *ActivityBuilder) WithSuccess(actionType string) *ActivityBuilder {
	b.Description.Resource = "cisco.activity." + actionType + "." + activityResourceSuccess
	b.Severity.Event = SeverityInformational
	b.Severity.Display = DisplaySeverityGood
	return b
}

func (b *ActivityBuilder) WithFailure(actionType string) *ActivityBuilder {
	b.Description.Resource = "cisco.activity." + actionType + "." + activityResourceFailure
	b.Severity.Event = SeverityCritical
	b.Severity.Display = DisplaySeverityCritical
	return b
}

func (b *ActivityBuilder) WithAction(action string) *ActivityBuilder {
	b.Action = action
	return b
}

func (b *ActivityBuilder) WithSeverity(severity string) *ActivityBuilder {
	b.Severity.Event = severity
	return b
}

func (b *ActivityBuilder) WithDisplaySeverity(displaySeverity string) *ActivityBuilder {
	b.Severity.Display = displaySeverity
	return b
}

func (b *ActivityBuilder) WithUser(user ActivityEntityReference) *ActivityBuilder {
	b.Operands.User = user
	return b
}

func (b *ActivityBuilder) displaySeverity() string {
	if b.Severity.Display != "" {
		return b.Severity.Display
	}

	switch b.Severity.Event {
	case SeverityInformational:
		return DisplaySeverityGood
	case SeverityWarning:
		return DisplaySeverityPoor
	case SeverityCritical:
		return DisplaySeverityCritical
	default:
		return DisplaySeverityUnknown
	}
}

func (b *ActivityBuilder) Message(ctx context.Context) (Message, error) {
	return b.Build(ctx)
}

func (b *ActivityBuilder) Validate() error {
	return types.ErrorMap{
		"description.text":     validation.Validate(&b.Description.Text, validation.Required),
		"description.resource": validation.Validate(&b.Description.Resource),
		"description.template": validation.Validate(&b.Description.Template),
		"serviceType":          validation.Validate(&b.ServiceType, validation.Required),
		"action":               validation.Validate(&b.Action, validation.In(ActionCreate, ActionUpdate, ActionDelete, ActionForceDelete)),
		"severity.display":     validation.Validate(&b.Severity.Display),
		"severity.event":       validation.Validate(&b.Severity.Event, validation.Required),
		"operands.object.id":   validation.Validate(&b.Operands.Object.Id, validation.Required),
		"operands.object.name": validation.Validate(&b.Operands.Object.Name, validation.Required),
		"operands.object.type": validation.Validate(&b.Operands.Object.Type, validation.Required),
		"operands.target.id":   validation.Validate(&b.Operands.Target.Id),
		"operands.target.name": validation.Validate(&b.Operands.Target.Name),
		"operands.target.type": validation.Validate(&b.Operands.Target.Type),
		"operands.tenant.id":   validation.Validate(&b.Operands.Tenant.Id, validation.Required),
		"operands.tenant.name": validation.Validate(&b.Operands.Tenant.Name),
	}
}

func (b *ActivityBuilder) Build(ctx context.Context) (msg Message, err error) {
	msg, err = NewMessage(ctx)
	if err != nil {
		return
	}

	svc, _ := config.FromContext(ctx).String("info.app.name")
	if svc == "" {
		svc = b.ServiceType
	}

	msg.Subtype = SubTypeActivity
	msg.Service = svc
	msg.Severity = b.displaySeverity()
	msg.AddKeyword(msg.Severity)

	if b.Operands.Tenant.Name != "" {
		msg.Security.TenantName = b.Operands.Tenant.Name
	}
	msg.Security.TenantId = b.Operands.Tenant.Id
	msg.AddKeyword(msg.Security.TenantId)

	if b.Operands.User.Name != "" {
		msg.Security.Username = b.Operands.User.Name
		msg.Security.OriginalUsername = b.Operands.User.Name
	}
	if b.Operands.User.Id != "" {
		msg.Security.UserId = b.Operands.User.Id
	}

	if b.Action != "" {
		msg.Action = b.Action
		msg.AddDetailWithKeyword(DetailsAction, b.Action)
	}

	msg.AddDetailWithKeyword(DetailsSeverity, b.Severity.Event)
	msg.AddDetailWithKeyword(DetailsServiceType, b.ServiceType)
	msg.AddDetail(DetailsDescriptionResource, b.Description.Resource)

	msg.AddDetail(DetailsObjectName, b.Operands.Object.Name)
	msg.AddDetailWithKeyword(DetailsObjectId, b.Operands.Object.Id)
	msg.AddDetailWithKeyword(DetailsObjectType, b.Operands.Object.Type)

	msg.AddDetail(DetailsTargetName, b.Operands.Target.Name)
	msg.AddDetail(DetailsTargetId, b.Operands.Target.Id)
	msg.AddDetail(DetailsTargetType, b.Operands.Target.Type)

	if b.Description.Text == "" {
		// Load the resource string if we have no template
		if b.Description.Template == "" && b.Description.Resource != "" && b.Description.ResourceLoader != nil {
			loader := b.Description.ResourceLoader
			b.Description.Template, err = loader(ctx, b.Description.Resource)
			if err != nil {
				return msg, err
			}
		}

		// Render the template if we have a template
		if b.Description.Template != "" {
			t := template.DollarTemplate(b.Description.Template)
			o := template.DollarRenderOptions{
				Variables: msg.Details.Strings(),
			}
			b.Description.Text = t.RenderString(o)
		}

	}
	msg.Description = b.Description.Text

	if err = validate.Validate(b); err != nil {
		return
	}

	return
}
