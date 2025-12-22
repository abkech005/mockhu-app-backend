package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// SESEmailService implements auth.EmailService using AWS SES v2
type SESEmailService struct {
	client      *sesv2.Client
	senderEmail string
}

// NewSESEmailService creates a new SES email service
func NewSESEmailService(ctx context.Context, region, senderEmail string) (*SESEmailService, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := sesv2.NewFromConfig(cfg)

	return &SESEmailService{
		client:      client,
		senderEmail: senderEmail,
	}, nil
}

// SendVerificationEmail sends a verification code to the user
func (s *SESEmailService) SendVerificationEmail(ctx context.Context, toEmail, code string) error {
	subject := "Verify your email - Mockhu"
	body := fmt.Sprintf(`
		<h1>Welcome to Mockhu!</h1>
		<p>Please use the following code to verify your email address:</p>
		<h2>%s</h2>
		<p>This code will expire in 10 minutes.</p>
		<p>If you did not request this, please ignore this email.</p>
	`, code)

	return s.sendEmail(ctx, toEmail, subject, body)
}

// SendPasswordResetEmail sends a password reset code
func (s *SESEmailService) SendPasswordResetEmail(ctx context.Context, toEmail, code string) error {
	subject := "Reset your password - Mockhu"
	body := fmt.Sprintf(`
		<h1>Password Reset Request</h1>
		<p>You requested to reset your password. Use the code below:</p>
		<h2>%s</h2>
		<p>This code will expire in 10 minutes.</p>
		<p>If you did not request this, please ignore this email.</p>
	`, code)

	return s.sendEmail(ctx, toEmail, subject, body)
}

func (s *SESEmailService) sendEmail(ctx context.Context, toEmail, subject, htmlBody string) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.senderEmail),
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(htmlBody),
						Charset: aws.String("UTF-8"),
					},
					Text: &types.Content{
						Data:    aws.String("Please enable HTML to view this email."), // Fallback
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	return nil
}
