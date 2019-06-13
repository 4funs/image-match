# Image match

## Sample using
```golang
tf, err := os.Open("image/wechat_moment_group.jpg")
checkerr(err)
defer tf.Close()
template, err := NewTemplateFromStream(tf, 1)
checkerr(err)
defer template.Close()

f, err := os.Open("image/wechat_moment.jpg")
checkerr(err)
defer f.Close()

target, err := NewMatFromStream(f)
checkerr(err)
defer target.Close()
result, err := template.Match(tt.args.img)
checkerr(err)
if result {
	// Matched
}
```