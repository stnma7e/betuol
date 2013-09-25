#version 330

smooth in vec4 interpColor;
smooth out vec4 outputColor;

void main()
{
	outputColor = interpColor;
}
