package imagematch

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	"gocv.io/x/gocv"
)

func TestImageMatch(t *testing.T) {
	type args struct {
		template string
		target   string
		sill     float32
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				template: "image/wechat_moment_group.jpg",
				target:   "image/wechat_moment.jpg",
				sill:     1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "reverse",
			args: args{
				target:   "image/wechat_moment_group.jpg",
				template: "image/wechat_moment.jpg",
				sill:     1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				template: "image/wechat_moment.jpg",
				sill:     1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "not exist 2",
			args: args{
				target: "image/wechat_moment.jpg",
				sill:   1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "not match",
			args: args{
				template: "image/wechat_moment.jpg",
				target:   "image/golang.jpeg",
				sill:     1,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImageMatch(tt.args.template, tt.args.target, tt.args.sill)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImageMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplate_Match(t *testing.T) {
	tf, err := os.Open("image/wechat_moment_group.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()
	template, err := NewTemplateFromStream(tf, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer template.Close()

	f, err := os.Open("image/wechat_moment.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	target, err := NewMatFromStream(f)
	if err != nil {
		t.Fatal(err)
	}
	defer target.Close()
	// fmt.Println("template:", target.Channels(), target.Empty(), target.Type())
	// fmt.Println("target:", target.Channels(), target.Empty(), target.Type())
	type args struct {
		img gocv.Mat
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				img: target,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.Match(tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Template.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
